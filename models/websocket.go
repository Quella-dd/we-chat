package models

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type WebsocketManager struct {
	Connects map[string]*websocket.Conn
}

func NewWebSocketManager() *WebsocketManager {
	return &WebsocketManager{
		Connects: make(map[string]*websocket.Conn),
	}
}

func (wm *WebsocketManager) Handler(ctx *gin.Context, userID string) {
	if ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil); err != nil {
		fmt.Println("websocket creat failed", err)
	} else {
		wm.Connects[userID] = ws
		if err := ws.WriteJSON("connect ws successd"); err != nil {
			fmt.Println("send Message errorr")
		}
	}
}
