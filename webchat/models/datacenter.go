package models

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type DataManager struct{}

type Message struct {
	Source      string `form:"source"`
	RoomID      *string
	Destination string `form:"destination"`
	Body        string `form:"body"`
}

func NewDataManager() *DataManager {
	return &DataManager{}
}

func (*DataManager) HandlerMessage(ctx *gin.Context) {
	if msg, err := parseMessage(ctx); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg)

		if msg.RoomID == nil {
			if ws, ok := ManageEnv.WebsocketManager.Connes[msg.Destination]; ok {
				ws.WriteJSON(msg)
			}
		}
	}
}

func getWsConn(dest string) {

}

func parseMessage(ctx *gin.Context) (*Message, error) {
	var msg Message
	if err := ctx.ShouldBind(&msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
