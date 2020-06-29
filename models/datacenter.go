package models

import (
	"database/sql/driver"
	"encoding/json"
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

type MessageRequest struct {
	SourceID int `form:"SourceID"`
	RoomID   *string
	DestID   int    `form:"DestID"`
	Body     string `form:"body"`
}

type MessageResponse struct {
	gorm.Model

	SourceID int
	RoomID   *string
	DestID   int
	Body     MessagesBody ` sql:"TYPE:json"`
}

type MessagesBody []MesageBody

type MesageBody struct {
	Content string
}

func (msg MessagesBody) Value() (driver.Value, error) {
	b, err := json.Marshal(msg)
	return string(b), err
}

func (msg *MessagesBody) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), msg)
}

func NewDataManager() *DataManager {
	database.DB.AutoMigrate(&MessageResponse{})
	return &DataManager{
		Publisher: pubsub.NewPublisher(timeout, buffer),
	}
}

// HandlerMessage save and send
func (dataManager *DataManager) HandlerMessage(ctx *gin.Context, userID string) {
	var msg MessageRequest
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

func (dataManager *DataManager) SaveMessage(msg *MessageRequest) error {
	var message MessageResponse
	if err := database.DB.Where("source_id = ? and dest_id = ?", msg.SourceID, msg.DestID).Find(&message).Error; err != nil {
		message.SourceID = msg.SourceID
		message.DestID = msg.DestID
		message.Body = []MesageBody{
			MesageBody{Content: msg.Body},
		}

		if err := database.DB.Create(&message).Error; err != nil {
			return err
		}
	} else {
		message.Body = append(message.Body, MesageBody{
			Content: msg.Body,
		})

		if err := database.DB.Model(&message).Update("body", message.Body).Error; err != nil {
			return err
		}
	}
	return nil
}

func (*DataManager) GetMessage(ctx *gin.Context, userID, destID string) error {
	var messages []*MessageResponse

	if err := database.DB.Where("source_id = ? and dest_id = ?", userID, destID).Find(&messages).Error; err != nil {
		return err
	}
	common.HttpSuccessResponse(ctx, messages)
	return nil
}
