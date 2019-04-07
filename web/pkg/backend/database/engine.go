package database

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/analogj/lantern/web/pkg/backend"
	"github.com/analogj/lantern/common/models"
	"github.com/chromedp/cdproto"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"time"
	"github.com/chromedp/cdproto/network"
	wsWrapper "github.com/analogj/lantern/web/pkg/models"
)

func New(toFrontend *chan wsWrapper.Wrapper, toBackend *chan wsWrapper.Wrapper, connString string) backend.Interface {
	e := new(engine)
	e.toFrontend = *toFrontend
	e.toBackend = *toBackend
	e.connString = connString
	return e
}

type engine struct {
	toFrontend chan<- wsWrapper.Wrapper // (send-only) send messages/requests on this channel to frontend clients.
	toBackend  <-chan wsWrapper.Wrapper // (listen-only) listen to this channel to retrieve messages from backend.
	connString string
}

// listen to backend messages from Database (Postgres)
func (e *engine) ListenMessages() {

	//open database connection
	orm, err := gorm.Open("postgres", e.connString)
	if err != nil {
		panic(err)
	}
	defer orm.Close()

	//setup raw sql connection to postgres for event notifications.
	_, err = sql.Open("postgres", e.connString)
	if err != nil {
		panic(err)
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	listener := pq.NewListener(e.connString, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen("events")
	if err != nil {
		panic(err)
	}

	fmt.Println("Start listening to messages in backend engine...")

	for {
		select {
		case postgresMsg := <-listener.Notify: //blocking
			fmt.Println("received postgres message", postgresMsg)
			//toFrontend messages should be immediately broadcasted

			fmt.Println("Received data from channel [", postgresMsg.Channel, "] :")
			logPayload(postgresMsg.Extra)

			event := models.DbNotify{}
			err = json.Unmarshal([]byte(postgresMsg.Extra), &event)
			if err != nil {
				fmt.Println("Error unmarshalling JSON: ", err)
				return
			}

			wsEvent, dbModel, err := processDatabaseEvent(orm, event.Table, event.Id)
			if err != nil {
				fmt.Println("Error transforming db event to ws event:: ", err)
				continue
			}

			// broadcast this event to the websocket.
			fmt.Println("forwarding parsed event to frontend...")
			e.toFrontend <- wsWrapper.Wrapper{Message: wsEvent}

			if wsEvent.Method == cdproto.EventNetworkResponseReceived {
				//this is a response, we need to send a couple of additional messages to "finalize" the response.
				dbResponse, ok := dbModel.(models.DbResponse)

				if !ok {
					//not a DB response, not sure what this is, bailing out.
					continue
				}

				eventDataReceived, err := generateDataReceivedEvent(dbResponse)
				if err != nil {
					continue
				}

				eventLoadingFinished, err := generateLoadingFinishedEvent(dbResponse)
				if err != nil {
					continue
				}

				e.toFrontend <- wsWrapper.Wrapper{Message: eventDataReceived}
				e.toFrontend <- wsWrapper.Wrapper{Message: eventLoadingFinished}

			}

		case frontendWrapper := <-e.toBackend:
			//Message requests forwarded from frontend.

			fmt.Println("received websocket command from frontend", frontendWrapper)

			switch frontendWrapper.Message.Method {
			case cdproto.CommandNetworkEnable:
				//do a database query for requests and responses to backfill the frontent.
				requests := []models.DbRequest{}
				orm.Find(&requests)

				for _, dbReq := range requests {
					message,  err := generateRequestWillBeSentEvent(dbReq)
					if err != nil {
						continue
					}

					e.toFrontend <- wsWrapper.Wrapper{Message: message, Destination: frontendWrapper.Destination}

					//find associated response
					dbResp := models.DbResponse{}
					orm.First(&dbResp, "request_id = ?", dbReq.Id)


					eventNetworkResponseReceived, err := generateNetworkResponseReceived(dbReq, dbResp)
					if err != nil {
						continue
					}

					eventDataReceived, err := generateDataReceivedEvent(dbResp)
					if err != nil {
						continue
					}

					eventLoadingFinished, err := generateLoadingFinishedEvent(dbResp)
					if err != nil {
						continue
					}

					e.toFrontend <- wsWrapper.Wrapper{Message: eventNetworkResponseReceived, Destination: frontendWrapper.Destination}
					e.toFrontend <- wsWrapper.Wrapper{Message: eventDataReceived, Destination: frontendWrapper.Destination}
					e.toFrontend <- wsWrapper.Wrapper{Message: eventLoadingFinished, Destination: frontendWrapper.Destination}
				}





			case cdproto.CommandNetworkGetResponseBody:
				params := network.GetResponseBodyParams{}
				if err := json.Unmarshal(frontendWrapper.Message.Params, &params); err != nil {
					//TODO:log an error message
					fmt.Println("An error occured parsing response body request params")
					continue
				}
				dbresp := models.DbResponse{}
				orm.First(&dbresp, "request_id = ?", params.RequestID)

				fmt.Println("Found ")
				payload := network.GetResponseBodyReturns{
					Body:          dbresp.Body,
					Base64encoded: true,
				}

				if payloadJson, err := payload.MarshalJSON(); err == nil {
					wsresp := cdproto.Message{
						ID:     frontendWrapper.Message.ID,
						Result: payloadJson,
					}
					e.toFrontend <- wsWrapper.Wrapper{Destination: frontendWrapper.Destination, Message: wsresp}
				}
			case cdproto.CommandNetworkGetRequestPostData:

				params := network.GetRequestPostDataParams{}
				if err := json.Unmarshal(frontendWrapper.Message.Params, &params); err != nil {
					//TODO:log an error message
					fmt.Println("An error occured parsing request post data request params")
					continue
				}
				dbreq := models.DbRequest{}
				orm.First(&dbreq, "id = ?", params.RequestID)

				fmt.Println("Found ")
				payload := network.GetRequestPostDataReturns{
					PostData:          dbreq.Body,
				}

				if payloadJson, err := payload.MarshalJSON(); err == nil {
					wsresp := cdproto.Message{
						ID:     frontendWrapper.Message.ID,
						Result: payloadJson,
					}
					e.toFrontend <- wsWrapper.Wrapper{Destination: frontendWrapper.Destination, Message: wsresp}
				}
			}

		}
	}
}

func logPayload(payload string) {
	// Prepare notification payload for pretty print
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(payload), "", "\t")
	if err != nil {
		fmt.Println("Error processing JSON: ", err)
		return
	}
	fmt.Println(string(prettyJSON.Bytes()))
}
