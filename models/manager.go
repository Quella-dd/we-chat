package models

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
	"webchat/common"
	"webchat/database"

	"github.com/docker/docker/pkg/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

const (
	UserActive      = "user_active"
)

var (
	buffer      = 1024
	timeout     = time.Second * 5
	onlineCount = 100
)

type DataCenterManager struct {
	Count         int
	TimeOut       time.Duration
	Redis         *redis.Client
	Publisher     *pubsub.Publisher
	ActiveChannel []chan string
}

func NewDataCenterManager() *DataCenterManager {
	database.DB.AutoMigrate(&SessionMessage{})
	dataCenterManager := &DataCenterManager{
		Count:     10,
		TimeOut:   time.Hour,
		Redis:     database.Redis,
		Publisher: pubsub.NewPublisher(timeout, buffer),
	}
	go dataCenterManager.InitPubsub()
	return dataCenterManager
}

func (dataCenter *DataCenterManager) InitPubsub() {
	pubsub := dataCenter.Redis.Subscribe(UserActive)

	for msg := range pubsub.Channel() {
		fmt.Println("ready to send message of redis")
		time.Sleep(time.Second)

		var redisSendSuccess bool = true

		for _, message := range dataCenter.Redis.LRange(msg.Payload, 0, -1).Val() {
			var redisMessage RequestBody
			if err := json.Unmarshal([]byte(message), &redisMessage); err != nil {
				fmt.Println("read redis error", err)
			} else {
				fmt.Printf("parse success: %+v\n", redisMessage)
				if err := ManageEnv.WebsocketManager.SendUserMessage(msg.Payload, redisMessage); err != nil {
					redisSendSuccess = false
				}
			}
		}

		if redisSendSuccess {
			dataCenter.Redis.Del(msg.Payload)
			fmt.Println("redis message send to user success")
		}
	}
}

func (dataCenter *DataCenterManager) HandlerMessage(ctx *gin.Context, userID string) error {
	var msg RequestBody
	if err := ctx.ShouldBind(&msg); err != nil {
		fmt.Println("parse error", err)
	}

	msg.CreateAt = time.Now().Unix()
	if err := dataCenter.Distribution(msg); err != nil {
		return err
	}
	if err := dataCenter.Save(msg); err != nil {
		return err
	}
	return nil
}

func (dataCenter *DataCenterManager) Distribution(msg RequestBody) error {
	var err error
	switch msg.Type {
	case common.USERMESSAGE:
		err = ManageEnv.WebsocketManager.SendUserMessage(msg.getUserIdentify(), msg)
		break
	case common.ROOMMESSAGE:
		err = ManageEnv.WebsocketManager.SendRoomMessage(msg)
		break
	case common.BORDERCASTMESSAGE:
		err = ManageEnv.WebsocketManager.SendBordcastMessage(msg)
		break
	}
	return err
}

func (dataCenter *DataCenterManager) Save(msg RequestBody) error {
	var message SessionMessage

	if err := database.DB.Where("source_id = ? and destination_id = ?", msg.SourceID, msg.DestinationID).Find(&message).Error; err != nil {
		message.SourceID = msg.SourceID
		message.DestinationID = msg.DestinationID
		message.Messages = append(message.Messages, MessagesBody{
			Create_At: msg.CreateAt,
			Content:   msg.Content,
		})

		if err := database.DB.Create(&message).Error; err != nil {
			return err
		}
	}

	message.Messages = append(message.Messages, MessagesBody{
		Create_At: msg.CreateAt,
		Content:   msg.Content,
	})
	database.DB.Model(&message).Update("messages", message.Messages)
	return nil
}

func (*DataCenterManager) GetMessage(ctx *gin.Context, userID, destID string) error {
	var resultMessages MessageInfos

	var requestMessage, reponseMessage SessionMessage

	if err := database.DB.Where("source_id = ? and destination_id = ?", userID, destID).First(&requestMessage).Error; err == nil {
		resultMessages = append(resultMessages, *(requestMessage.GetMessageInfo())...)
	}

	if err := database.DB.Where("source_id = ? and destination_id = ?", destID, userID).First(&reponseMessage).Error; err == nil {
		resultMessages = append(resultMessages, *(reponseMessage.GetMessageInfo())...)
	}

	sort.Sort(resultMessages)
	common.HttpSuccessResponse(ctx, resultMessages)
	return nil
}
