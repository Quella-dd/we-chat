package models

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"we-chat/message"
	Message "we-chat/message"

	"github.com/docker/docker/pkg/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

const (
	UserActive = "user_active"
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
	dataCenterManager := &DataCenterManager{
		Count:     10,
		TimeOut:   time.Hour,
		Redis:     ManagerEnv.redis,
		Publisher: pubsub.NewPublisher(timeout, buffer),
	}
	go dataCenterManager.InitPubsub()
	return dataCenterManager
}

func (dataCenter *DataCenterManager) InitPubsub() {
	pubsub := dataCenter.Redis.Subscribe(UserActive)

	for msg := range pubsub.Channel() {
		time.Sleep(time.Second)
		var redisSendSuccess bool = true
		for _, msg := range dataCenter.Redis.LRange(msg.Payload, 0, -1).Val() {
			var redisMessage Message.RequestMessage
			if err := json.Unmarshal([]byte(msg), &redisMessage); err != nil {
				fmt.Println("read redis error", err)
			} else {
				fmt.Printf("parse success: %+v\n", redisMessage)
				if err := ManagerEnv.WebsocketManager.SendUserMessage(redisMessage, redisMessage.Scope.DestinationID); err != nil {
					redisSendSuccess = false
				}
			}
		}

		if redisSendSuccess {
			dataCenter.Redis.Del(msg.Payload)
			log.Fatal("redis message send to user success")
		}
	}
}

// 1. create or update session
// 2. Distribution =>  message.USERMESSAGE, message.ROOMMESSAGE, message.BORDERCASTMESSAGE
// 3. Save with Mysql
func (dataCenter *DataCenterManager) HandlerMessage(c *gin.Context, requestMessage message.RequestMessage) error {
	var session *Session

	sessionID, _ := requestMessage.Scope.GetSession()
	if err := ManagerEnv.DB.Find("id = ?", sessionID).Find(session).Error; err != nil {
		session, err = ManagerEnv.SessionManager.CreateSession(&Session{
			Owner:         requestMessage.Scope.SourceID,
			Src:           requestMessage.Scope.DestinationID,
			LatestTime:    time.Now(),
			LatestContent: requestMessage.Content,
		})
		if err != nil {
			return fmt.Errorf("create Session failed, error: %+v\n", err)
		}

		requestMessage.ID = session.ID

		if err := ManagerEnv.DB.Create(&requestMessage).Error; err != nil {
			return err
		}
	}

	// requestMessage.SessionID = sessionID
	if err := dataCenter.Distribution(requestMessage); err != nil {
		return err
	}

	return dataCenter.Save(requestMessage, *session)
}

func (dataCenter *DataCenterManager) Distribution(msg Message.RequestMessage) error {
	var err error
	switch msg.Scope.Stype {
	case message.USERMESSAGE:
		err = ManagerEnv.WebsocketManager.SendUserMessage(msg, msg.Scope.DestinationID)
		break
	case message.ROOMMESSAGE:
		err = ManagerEnv.WebsocketManager.SendRoomMessage(msg)
		break
	case message.BORDERCASTMESSAGE:
		err = ManagerEnv.WebsocketManager.SendBordcastMessage(msg)
		break
	}
	return err
}

// 离线消息存储
func (dataCenter *DataCenterManager) Save(message Message.RequestMessage, session Session) error {
	if err := ManagerEnv.DB.Create(&message).Error; err != nil {
		return err
	}

	session.LatestContent = message.Content
	if err := ManagerEnv.DB.Update("LatestContent", session.LatestContent).Error; err != nil {
		return err
	}
	return nil
}
