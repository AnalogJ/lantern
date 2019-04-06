package main

import (
	"net/http"

	db "github.com/analogj/lantern/web/pkg/backend/database"
	ws "github.com/analogj/lantern/web/pkg/frontend/websocket"
	"log"
	"github.com/analogj/lantern/web/pkg/models"
)

var toFrontend = make(chan models.Wrapper ) // send to websocket (frontend)
var toBackend = make(chan models.Wrapper )  // send to database (backend)

func main() {

	//Initialize the frontend engine
	frontendEngine := ws.New(&toFrontend, &toBackend)

	//Initialize the backend engine
	backendEngine := db.New(&toFrontend, &toBackend, "host=database sslmode=disable dbname=lantern user=lantern password=lantern-password")

	// Create a simple file server
	http.Handle("/", http.FileServer(http.Dir("/srv/lantern/")))
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