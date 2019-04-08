package main

import (
	"net/http"

	db "github.com/analogj/lantern/web/pkg/backend/database"
	ws "github.com/analogj/lantern/web/pkg/frontend/websocket"
	"log"
	"github.com/analogj/lantern/web/pkg/models"
	"io/ioutil"
	"fmt"
	"strconv"
	"time"
	"bytes"
	"path/filepath"
)

var toFrontend = make(chan models.Wrapper ) // send to websocket (frontend)
var toBackend = make(chan models.Wrapper )  // send to database (backend)

func main() {

	//Initialize the frontend engine
	frontendEngine := ws.New(&toFrontend, &toBackend)

	//Initialize the backend engine
	backendEngine := db.New(&toFrontend, &toBackend, "host=database sslmode=disable dbname=lantern user=lantern password=lantern-password")

	// Create a simple file server
	http.HandleFunc("/certs/ca.crt", DownloadHandler)
	http.HandleFunc("/certs/lantern.mobileconfig", DownloadHandler)

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


func DownloadHandler(w http.ResponseWriter, r *http.Request) {

	file := filepath.Join("/srv/lantern/", r.RequestURI)

	downloadBytes, err := ioutil.ReadFile(file)

	if err != nil {
		fmt.Println(err)
	}

	// set the default MIME type to send
	mime := http.DetectContentType(downloadBytes)

	fileSize := len(string(downloadBytes))

	// Generate the server headers
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(file)+"")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Length", strconv.Itoa(fileSize))
	w.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")

	//b := bytes.NewBuffer(downloadBytes)
	//if _, err := b.WriteTo(w); err != nil {
	//              fmt.Fprintf(w, "%s", err)
	//      }

	// force it down the client's.....
	http.ServeContent(w, r, file, time.Now(), bytes.NewReader(downloadBytes))

}