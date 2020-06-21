package models

import (
	"fmt"
	"time"

	"github.com/docker/docker/pkg/pubsub"
	"github.com/gin-gonic/gin"
)

var (
	buffer  = 1024
	timeout = time.Second * 5
)

type DataManager struct {
	Publisher *pubsub.Publisher
}

type Message struct {
	Source      string `form:"source"`
	RoomID      *string
	Destination string `form:"destination"`
	Body        string `form:"body"`
}

type TestMessage struct {
	SourceName string `form:"SourceName"`
	Content    string `form:"content"`
}

func NewDataManager() *DataManager {
	return &DataManager{
		Publisher: pubsub.NewPublisher(timeout, buffer),
	}
}

func (dm *DataManager) HandlerMessage(ctx *gin.Context) {
	var msg TestMessage
	ctx.Request.ParseForm()

	if err := ctx.ShouldBind(&msg); err != nil {
		fmt.Println("parse error", err)
	} else {
		if ws, ok := ManageEnv.WebsocketManager.Connects[msg.SourceName]; ok {
			ws.WriteJSON(msg.Content)
		} else {
			fmt.Println("not found ws")
		}
	}

	// if msg, err := parseMessage(ctx); err != nil {
	// 	dm.Publisher.Publish(msg)
	// }

	// if msg, err := parseMessage(ctx); err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(msg)

	// 	if msg.RoomID == nil {
	// 		if ws, ok := ManageEnv.WebsocketManager.Connes[msg.Destination]; ok {
	// 			ws.WriteJSON(msg)
	// 		}
	// 	}
	// }
}

func parseMessage(ctx *gin.Context) (*Message, error) {
	var msg Message
	if err := ctx.ShouldBind(&msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
