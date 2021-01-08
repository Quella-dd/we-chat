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
	if err := ManageEnv.DataCenterManager.Redis.Publish(UserActive, id).Err(); err != nil {
		return err
	}

	return ws.WriteJSON(Event{
		Action: Login_Event,
	});
}

func (vm *WebsocketManager) SendUserMessage(identify string, msg message.RequestMessage) error {
	if ws, ok := vm.Connections.Load(identify); ok {
		if conn, ok := ws.(*websocket.Conn); ok {
			return conn.WriteJSON(msg)
		}
	}

	// fmt.Println("Destination user is offline")
	// logrus.Info("Destination user is offline")

	// 如果用户离线，将message保存到离线数据库， redis列表的Key 为identify (list)
	if err := ManageEnv.DataCenterManager.Redis.RPush(identify, msg).Err(); err != nil {
		return err
	}

	// value of return descide save the offline data in sql
	return nil
	// return errors.New(fmt.Sprintf("websocket conn recode not fond: %+v\n", identify))
}

func (vm *WebsocketManager) SendRoomMessage(msg message.RequestMessage) error {
	roomID := strconv.Itoa(msg.RoomID)
	room, err := ManageEnv.GroupManager.GetGroup(roomID)
	if err != nil {
		return err
	}
	for _, v := range room.Childes {
		vm.SendUserMessage(strconv.Itoa(int(v.ID)), msg)
	}
	return nil
}

func (vm *WebsocketManager) SendBordcastMessage(msg message.RequestMessage) error {
	users, err := ManageEnv.UserManager.ListUsers()
	if err != nil {
		return err
	}
	for _, v := range users {
		vm.SendUserMessage(strconv.Itoa(int(v.ID)), msg)
	}
	return nil
}
