package models

import (
	"net/http"
	"strconv"
	"sync"
	"we-chat/message"

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
	Connections *sync.Map
}

func NewWebSocketManager() *WebsocketManager {
	return &WebsocketManager{
		Connections: &sync.Map{},
	}
}

func (m *WebsocketManager) Handler(ctx *gin.Context, id string) error {
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return err
	}

	m.Connections.Store(id, ws)

	// publish event of userActive
	if err := ManagerEnv.DataCenterManager.Redis.Publish(UserActive, id).Err(); err != nil {
		return err
	}

	return ws.WriteJSON(Event{
		Action: Login_Event,
	});
}

func (m *WebsocketManager) SendUserMessage(msg message.RequestMessage, destinationID string) error {
	if ws, ok := m.Connections.Load(destinationID); ok {
		if conn, ok := ws.(*websocket.Conn); ok {
			return conn.WriteJSON(msg)
		}
	}

	// 如果用户离线，将message保存到离线数据库， redis列表的Key 为identify (list)
	if err := ManagerEnv.DataCenterManager.Redis.RPush(destinationID, msg).Err(); err != nil {
		return err
	}

	// value of return descide save the offline data in sql
	return nil
	// return errors.New(fmt.Sprintf("websocket conn recode not fond: %+v\n", identify))
}

func (m *WebsocketManager) SendRoomMessage(msg message.RequestMessage) error {
	room, err := ManagerEnv.GroupManager.GetGroup(msg.Scope.RoomID)
	if err != nil {
		return err
	}
	for _, v := range room.Childes {
		m.SendUserMessage(msg, v)
	}
	return nil
}

func (m *WebsocketManager) SendBordcastMessage(msg message.RequestMessage) error {
	users, err := ManagerEnv.UserManager.ListUsers()
	if err != nil {
		return err
	}
	for _, user := range users {
		userIDString := strconv.Itoa(int(user.ID))
		m.SendUserMessage(msg, userIDString)
	}
	return nil
}
