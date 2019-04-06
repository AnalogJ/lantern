package websocket

import (
	"fmt"
	"github.com/analogj/lantern/web/pkg/frontend"
	"github.com/chromedp/cdproto"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"github.com/analogj/lantern/web/pkg/models"
)

// Class constructor.

func New(toFrontend *chan models.Wrapper, toBackend *chan models.Wrapper) frontend.Interface {
	e := new(engine)
	e.clients = make(map[*websocket.Conn]bool)
	e.toFrontend = *toFrontend
	e.toBackend = *toBackend
	return e
}

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type engine struct {
	toFrontend chan models.Wrapper     // (listen & send) listen to this channel for messages to send to connected clients, also directly send responses to this channel
	toBackend  chan<- models.Wrapper   // (send-only) this is a channel that can be used to send messages/requests to the backend.
	clients    map[*websocket.Conn]bool // map of long lived connected clients
}

func (e *engine) RegisterConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	e.clients[ws] = true

	for {
		_, data, err := ws.ReadMessage() //blocking, wait for new message.
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		// Read in a new message as JSON and map it to a Message object
		wsCommand := cdproto.Message{}
		wsCommand.UnmarshalJSON(data)

		log.Printf("recieved msg: %v %s %v", wsCommand.ID, wsCommand.Method.String(), string(wsCommand.Params))

		if err != nil {
			log.Printf("error: %v", err)
			delete(e.clients, ws)
			break
		}
		// Send the newly received message to the toFrontend channel
		response := cdproto.Message{
			ID: wsCommand.ID,
		}

		switch wsCommand.Method.String() {

		//enable specific tabs
		case cdproto.CommandNetworkEnable:
			response.Result = []byte("{}")
			e.toFrontend <- models.Wrapper{ Message: response }
			//specifically forward the NetworkEnable command to the backend so that we can trigger a query of the Database
			// to get existing recordings
			e.toBackend <- models.Wrapper{Message: wsCommand, Destination: ws }

		//disable specific features
		case cdproto.CommandPageEnable,
			cdproto.CommandDOMEnable,
			cdproto.CommandRuntimeEnable,
			cdproto.CommandNetworkEmulateNetworkConditions,
			cdproto.CommandEmulationCanEmulate:
			response.Result = []byte(`{"result":false}`)
			e.toFrontend <- models.Wrapper{Message: response}

		//forward some commands to the backend (queries, etc)
		case cdproto.CommandNetworkGetResponseBody:
			e.toBackend <- models.Wrapper{Message: wsCommand, Destination: ws }

		//Fallback, say that we don't support this command.
		default:
			respErr := new(cdproto.Error)
			respErr.Message = fmt.Sprintf("Unsupported command: %v", wsCommand.Method.String())
			e.toFrontend <- models.Wrapper{Message: cdproto.Message{
				ID:    wsCommand.ID,
				Error: respErr,
			}}
		}
	}
}

func (e *engine) ListenMessages() {
	for {
		// Grab the next message from the toFrontend channel
		wrapper := <-e.toFrontend

		message := wrapper.Message
		destination := wrapper.Destination

		if wrapper.Destination != nil {

			log.Printf("sending frontend wrapper sent to one client: %v %s %v", message.ID, message.Method.String(), string(message.Result))

			// send the message to only one destination
			client := wrapper.Destination
			err := destination.WriteJSON(wrapper.Message)

			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(e.clients, destination)
			}
		} else {
			//send message to all clients.
			// Send it out to every client that is currently connected

			log.Printf("sending frontend wrapper sent to clients: %v %s %v", message.ID, message.Method.String(), string(message.Result))
			for client := range e.clients {
				err := client.WriteJSON(wrapper.Message)

				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(e.clients, client)
				}
			}
		}



	}
}
