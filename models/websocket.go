package models

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

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

func (wm *WebsocketManager) Handler(ctx *gin.Context, userID string) {
	if ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil); err != nil {
		fmt.Println("websocket creat failed", err)
	} else {
		wm.Connections.Store(userID, ws)
		if err := ws.WriteJSON("connect ws successd"); err != nil {
			fmt.Println("send Message errorr")
		}
	}
}

func (vm *WebsocketManager) SendUserMessage(identify string, msg SessionMessage) error {
	if ws, ok := vm.Connections.Load(identify); ok {
		if conn, ok := ws.(*websocket.Conn); ok {
			return conn.WriteJSON(msg)
		}
	}

	// TOOD: 此时应该将message保存到离线数据库， 当用户上线时进行数据推送
	return errors.New(fmt.Sprintf("websocket conn recode not fond: %+v\n", identify))
}

func (vm *WebsocketManager) SendRoomMessage(msg SessionMessage) error {
	roomID := strconv.Itoa(msg.RoomID)
	room, err := ManageEnv.RoomManager.GetRoom(roomID)
	if err != nil {
		return err
	}
	for _, v := range room.Childrens {
		vm.SendUserMessage(strconv.Itoa(int(v.ID)), msg)
	}
	return nil
}

func (vm *WebsocketManager) SendBordcastMessage(msg SessionMessage) error {
	users, err := ManageEnv.UserManager.listUsers()
	if err != nil {
		return err
	}
	for _, v := range users {
		vm.SendUserMessage(strconv.Itoa(int(v.ID)), msg)
	}
	return nil
}

