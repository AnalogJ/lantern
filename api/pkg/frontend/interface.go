package frontend

import "net/http"

type Interface interface {

	RegisterConnection(w http.ResponseWriter, r *http.Request)
	ListenMessages()
}
