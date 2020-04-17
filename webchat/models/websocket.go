package models

import (
	"fmt"
	"net/http"
	"time"

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
	Connes map[string] *websocket.Conn
}

func NewWebSocketManager() *WebsocketManager {
	return &WebsocketManager{
		Connes: make(map[string] *websocket.Conn),
	}
}

func (wm *WebsocketManager) Handler(ctx *gin.Context, userID string) {
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("websocket creat failed", err)
	} else {
		wm.Connes[userID] = ws
		time.Sleep(5 * time.Second)
		ws.WriteJSON(struct {
			UserID string
			Message string
			Now string
		}{userID, "hello", time.Now().Format("01-02 15:03:04")})
	}
}