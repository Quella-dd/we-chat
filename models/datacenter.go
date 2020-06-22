package models

import (
	"fmt"
	"strconv"
	"time"
	"webchat/common"
	"webchat/database"

	"github.com/docker/docker/pkg/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var (
	buffer  = 1024
	timeout = time.Second * 5
)

type DataManager struct {
	Publisher *pubsub.Publisher
}

type Message struct {
	gorm.Model

	SourceID int `form:"SourceID"`
	RoomID   *string
	DestID   int         `form:"DestID"`
	Body     string `form:"body"`
}

func NewDataManager() *DataManager {
	database.DB.AutoMigrate(&Message{})
	return &DataManager{
		Publisher: pubsub.NewPublisher(timeout, buffer),
	}
}

// HandlerMessage save and send
func (dataManager *DataManager) HandlerMessage(ctx *gin.Context, userID string) {
	var msg Message
	if err := ctx.ShouldBind(&msg); err != nil {
		fmt.Println("parse error", err)
	} else {
		// save
		msg.SourceID, _ = strconv.Atoi(userID)
		if err := dataManager.SaveMessage(&msg); err != nil {
			fmt.Println("save message error:", err)
		}

		// send
		if ws, ok := ManageEnv.WebsocketManager.Connects[strconv.Itoa(msg.SourceID)]; ok {
			ws.WriteJSON(msg.Body)
		} else {
			fmt.Println("not found ws")
		}
	}
}

func (dataManager *DataManager) SaveMessage(msg *Message) error {
	if err := database.DB.Create(msg).Error; err != nil {
		return err
	}
	return nil
}

func (*DataManager) GetMessage(ctx *gin.Context, userID string) error {
	var messages []*Message
	if err := database.DB.Find(&messages, userID).Error; err != nil {
		return err
	}
	common.HttpSuccessResponse(ctx, messages)
	return nil
}
