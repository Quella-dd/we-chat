package models

import (
	"github.com/gorilla/websocket"
)

type Listener struct {
	ws *websocket.Conn
}

func (l *Listener) Listen() {
	// ch := ManageEnv.DataManager.Publisher.SubscribeTopic()

	// for {
	// 	select {
	// 	case msg := <-ch:
	// 	l.ws.WriteJSON(msg)
	// 	}
	// }
}
