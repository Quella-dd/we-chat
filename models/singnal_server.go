package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"webchat/message"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var defaultExpiration = 10 * time.Second

func SendRTCRequest(c *gin.Context) error {
	var event message.RTCEvent
	if err := c.ShouldBind(&event).Error(); err != "" {
		return errors.New(fmt.Sprintf("bad request: %+v\n", http.StatusBadRequest))
	}

	event.ID = uuid.NewV4().String()

	// TODO: 当发起方发起的请求，接收方无任何响应，在固定的等待时间后，将cancel此次请求， 无论该事件成功与否
	// TODO: 都应该将此次事件序列化到数据库message中，type为RTC
	ManageEnv.DataCenterManager.Redis.Set(event.ID, event, defaultExpiration)

	// 发送方通知接收方接受RTC请求
	if ws, ok := ManageEnv.WebsocketManager.Connections.Load(event.DestinationID); !ok {
		fmt.Println("destination user offline")
	} else {
		if conn, ok := ws.(*websocket.Conn); ok {
			return conn.WriteJSON(event)
		}
	}
	return nil
}

// 等消息结束后，实例化消息在数据库中
func HandlerRTCRequest(c *gin.Context, eventID, status string) error {
	event, err := GetEventByID(eventID)
	if err != nil {
		return err
	}
	event.Status = status

	switch status {
	case "ack":
		event.StartTime = time.Now()
	case "refuse":
		ManageEnv.DataCenterManager.Redis.Del(eventID)
	}
	//// 接收方通知发送方拒绝或者接受RTC请求
	//if ws, ok := ManageEnv.WebsocketManager.Connections.Load(event.SourceID); !ok {
	//	fmt.Println("source user offline")
	//} else {
	//	if conn, ok := ws.(*websocket.Conn); ok {
	//		return conn.WriteJSON(event)
	//	}
	//}
	return nil
}

func HangDownRTCRequest(c *gin.Context, eventID string) error {
	return nil
}

func GetEventByID(eventID string) (*message.RTCEvent, error) {
	var event message.RTCEvent
	result := ManageEnv.DataCenterManager.Redis.Get(eventID).Val()
	if err := json.Unmarshal([]byte(result), &event); err != nil {
		return nil, err
	}
	return &event, nil
}
