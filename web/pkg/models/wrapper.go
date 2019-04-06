package models

import (
	"github.com/gorilla/websocket"
	"github.com/chromedp/cdproto"
)

type Wrapper struct {
	Destination *websocket.Conn
	Message cdproto.Message
}
