package database

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/analogj/lantern/api/pkg/backend"
	"github.com/analogj/lantern/common/models"
	"github.com/chromedp/cdproto"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"time"
)

func New(toFrontend *chan cdproto.Message, toBackend *chan cdproto.Message, connString string) backend.Interface {
	e := new(engine)
	e.toFrontend = *toFrontend
	e.toBackend = *toBackend
	e.connString = connString
	return e
}

type engine struct {
	toFrontend chan<- cdproto.Message // (send-only) send messages/requests on this channel to frontend clients.
	toBackend  <-chan cdproto.Message // (listen-only) listen to this channel to retrieve messages from backend.
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
			e.toFrontend <- wsEvent

			if wsEvent.Method == cdproto.EventNetworkResponseReceived {
				//this is a response, we need to send a couple of additional messages to "finalize" the response.
				dbResponse, ok := dbModel.(models.DbResponse)

				if !ok {
					//not a DB response, not sure what this is, bailing out.
					continue
				}

				eventDataReceived := cdproto.Message{
					Method: cdproto.EventNetworkDataReceived,
				}
				timestamp := cdp.MonotonicTime(dbResponse.CreatedAt)
				dataRecievedPayload := network.EventDataReceived{
					RequestID:         network.RequestID(fmt.Sprint(dbResponse.RequestId)),
					Timestamp:         &timestamp,
					DataLength:        dbResponse.ContentLength,
					EncodedDataLength: int64(len(dbResponse.Body)),
				}
				dataRecievedJsonBytes, err := json.Marshal(dataRecievedPayload)
				if err != nil {
					continue
				}
				eventDataReceived.Params = dataRecievedJsonBytes

				eventLoadingFinished := cdproto.Message{
					Method: cdproto.EventNetworkLoadingFinished,
				}
				loadingFinishedPayload := network.EventLoadingFinished{
					RequestID:         network.RequestID(fmt.Sprint(dbResponse.RequestId)),
					Timestamp:         &timestamp,
					EncodedDataLength: float64(len(dbResponse.Body)),
				}
				loadingFinishedJsonBytes, err := json.Marshal(loadingFinishedPayload)
				if err != nil {
					continue
				}
				eventLoadingFinished.Params = loadingFinishedJsonBytes

				e.toFrontend <- eventDataReceived
				e.toFrontend <- eventLoadingFinished

			}

		case frontendCmd := <-e.toBackend: //blocking
			fmt.Println("received websocket command from frontend", frontendCmd)

			switch frontendCmd.Method {
			case cdproto.CommandNetworkEnable:
				//TODO: do a database query for requests and responses and forward to frontend.

				//requests := []models.DbRequest{}
				//responses := []models.DbResponse{}
				//
				//

			case cdproto.CommandNetworkGetResponseBody:
				params := network.GetResponseBodyParams{}
				if err := json.Unmarshal(frontendCmd.Params, &params); err != nil {
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
						ID:     frontendCmd.ID,
						Result: payloadJson,
					}
					e.toFrontend <- wsresp
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
