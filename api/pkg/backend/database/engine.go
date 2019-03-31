package database

import (
	"github.com/analogj/lantern/api/pkg/backend"
	"github.com/chromedp/cdproto"
	"database/sql"
	"github.com/lib/pq"
	"fmt"
	"time"
	"bytes"
	"encoding/json"
	"github.com/analogj/lantern/common/models"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/cdp"
)

func New(toFrontend *chan cdproto.Message, toBackend *chan cdproto.Message, connString string) backend.Interface{
	e := new(engine)
	e.toFrontend = *toFrontend
	e.toBackend = *toBackend
	e.connString = connString
	return e
}

type engine struct {
	toFrontend chan<- cdproto.Message     // (send-only) send messages/requests on this channel to frontend clients.
	toBackend <-chan cdproto.Message     // (listen-only) listen to this channel to retrieve messages from backend.
	connString string
}


// listen to backend messages from Database (Postgres)
func (e *engine) ListenMessages(){

	_, err := sql.Open("postgres", e.connString)
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

			wsEvent, dbModel, err  := processDatabaseEvent(event.Table, event.Data)
			if err != nil {
				fmt.Println("Error transforming db event to ws event:: ", err)
				continue
			}

			// broadcast this event to the websocket.
			fmt.Println("forwarding parsed event to frontend...")
			e.toFrontend <- wsEvent

			if wsEvent.Method == cdproto.EventNetworkResponseReceived {
				//this is a response, we need to send a couple of additional messages to "finalize" the response.
				dbResponse, ok :=  dbModel.(models.DbResponse)

				if !ok {
					//not a DB response, not sure what this is, bailing out.
					continue
				}

				eventDataReceived := cdproto.Message{
					Method: cdproto.EventNetworkDataReceived,
				}
				timestamp := cdp.MonotonicTime(dbResponse.CreatedAt)
				dataRecievedPayload := network.EventDataReceived{
					RequestID: network.RequestID(fmt.Sprint(dbResponse.RequestId)),
					Timestamp: &timestamp,
					DataLength: -1,
					EncodedDataLength: -1,
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
					RequestID: network.RequestID(fmt.Sprint(dbResponse.RequestId)),
					Timestamp: &timestamp,
					EncodedDataLength: -1,
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

			case cdproto.CommandNetworkGetResponseBody:
				//TODO: do a database query for the response body and forward to frontend
				//toBackend messages should be parsed, and the
			}
		}
	}

}

func logPayload(payload string){
	// Prepare notification payload for pretty print
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(payload), "", "\t")
	if err != nil {
		fmt.Println("Error processing JSON: ", err)
		return
	}
	fmt.Println(string(prettyJSON.Bytes()))
}