package main

import (
	"net/http"

	"github.com/chromedp/cdproto"
	ws "github.com/analogj/lantern/api/pkg/frontend/websocket"
	db "github.com/analogj/lantern/api/pkg/backend/database"

	"log"
)


var toFrontend = make(chan cdproto.Message) // send to websocket (frontend)
var toBackend = make(chan cdproto.Message)  // send to database (backend)

func main() {


	//Initialize the frontend engine
	frontendEngine := ws.New(&toFrontend, &toBackend)

	//Initialize the backend engine
	backendEngine := db.New(&toFrontend, &toBackend, "host=database sslmode=disable dbname=lantern user=lantern password=lantern-password")


	// Create a simple file server
	http.Handle("/", http.FileServer(http.Dir("../public")))
	// Configure websocket route
	http.HandleFunc("/ws", frontendEngine.RegisterConnection)


	go frontendEngine.ListenMessages()
	go backendEngine.ListenMessages()

	// Start the server on localhost port 9000 and log any errors
	log.Println("http server started on :9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
