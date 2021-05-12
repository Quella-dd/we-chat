package models

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"we-chat/message"
	Message "we-chat/message"

	"github.com/docker/docker/pkg/pubsub"
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
				// fmt.Printf("parse success: %+v\n", redisMessage)
				// if err := ManagerEnv.WebsocketManager.SendUserMessage(redisMessage, redisMessage.Scope.DestinationID); err != nil {
				// 	redisSendSuccess = false
				// }
				var resultRequestMessage Message.RequestMessage

				resultRequestMessage.DestinationID = redisMessage.OwnerID
				resultRequestMessage.OwnerID = redisMessage.DestinationID
				resultRequestMessage.Content = redisMessage.Content

				session, err := ManagerEnv.DataCenterManager.UpdateOrCreateSession(resultRequestMessage)
				if err != nil {
					fmt.Println("Send UserMessage Error:", err)
				}
				redisMessage.SessionID = fmt.Sprintf("%+v", session.ID)

				if err := ManagerEnv.DataCenterManager.Save(redisMessage, *session); err != nil {
					fmt.Println("pubsub HandlerMessage Error:", err)
				}
			}
		}

		if redisSendSuccess {
			dataCenter.Redis.Del(msg.Payload)
			log.Fatal("redis message send to user success")
		}
	}
}

func (dataCenter *DataCenterManager) HandlerMessage(requestMessage message.RequestMessage) error {
	session, err := dataCenter.UpdateOrCreateSession(requestMessage)
	if err != nil {
		return err
	}

	requestMessage.SessionID = fmt.Sprintf("%+v", session.ID)

	user, err := ManagerEnv.UserManager.GetUser(requestMessage.OwnerID, "id")
	if err == nil {
		requestMessage.OwnerName = user.Name
	}

	// 消息存储在数据库
	if err := dataCenter.Save(requestMessage, *session); err != nil {
		return nil
	}

	// 消息下发
	if err := dataCenter.Distribution(requestMessage); err != nil {
		return err
	}
	return nil
}

func (dataCenter *DataCenterManager) UpdateOrCreateSession(requestMessage message.RequestMessage) (*Session, error) {
	var session *Session
	var err error

	session, err = getSessionwithScope(requestMessage.Scope)

	if err != nil {
		session, err = ManagerEnv.SessionManager.CreateSession(&Session{
			LatestTime: time.Now(),
			Scope: Message.Scope{
				Stype:         requestMessage.Scope.Stype,
				RoomID:        requestMessage.Scope.RoomID,
				OwnerID:       requestMessage.Scope.OwnerID,
				DestinationID: requestMessage.Scope.DestinationID,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("create Session failed, error: %+v\n", err)
		}
	}
	return session, nil
}

func (dataCenter *DataCenterManager) GetMessages(id string) ([]message.RequestMessage, error) {
	var messages []message.RequestMessage

	if err := ManagerEnv.DB.Where("session_id = ?", id).Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil

	// for _, message := range messages {
	// 	user, err := ManagerEnv.UserManager.GetUser(message.OwnerID, "id")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	} else {
	// 		message.OwnerName = user.Name
	// 	}
	// }
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

	// some bug
	if err := ManagerEnv.DB.Find(&session).Update("LatestContent", session.LatestContent).Error; err != nil {
		return err
	}
	return nil
}

func getSessionwithScope(scope Message.Scope) (*Session, error) {
	var session Session

	if scope.RoomID != "" {
		if err := ManagerEnv.DB.Where("owner_id = ? AND destination_id = ?", scope.OwnerID, scope.DestinationID).Find(&session).Error; err != nil {
			return nil, err
		}
	} else {
		if err := ManagerEnv.DB.Where("owner_id = ? AND destination_id = ? AND room_id = ?", scope.OwnerID, scope.DestinationID, session.RoomID).Find(&session).Error; err != nil {
			return nil, err
		}
	}
	return &session, nil
}
