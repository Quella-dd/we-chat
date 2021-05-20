package models

import (
	"fmt"
	"strconv"
	"we-chat/event"
	"we-chat/message"
	Message "we-chat/message"
)

// var (
// 	upgrader = websocket.Upgrader{
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true
// 		},
// 	}
// )

type WebsocketManager struct {
	// Connections *sync.Map
}

func NewWebSocketManager() *WebsocketManager {
	return &WebsocketManager{
		// Connections: &sync.Map{},
	}
}

// func (m *WebsocketManager) InitWs(c *gin.Context, id string) error {
// 	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		return err
// 	}

// 	m.Connections.Store(id, ws)

// 	if err := ManagerEnv.DataCenterManager.Redis.Publish(UserActive, id).Err(); err != nil {
// 		return err
// 	}
// 	return ws.WriteJSON(Event{
// 		Action: Login_Event,
// 	})
// }

type MessageEvt struct {
	Content interface{}
}

func (e MessageEvt) Type() string {
	return "message-evt"
}

func (m *WebsocketManager) SendUserMessage(msg message.RequestMessage, destinationID string) error {
	var resultRequestMessage Message.RequestMessage

	resultRequestMessage.Stype = msg.Stype
	resultRequestMessage.Content = msg.Content
	resultRequestMessage.OwnerID = destinationID

	if msg.RoomID != "" {
		resultRequestMessage.RoomID = msg.RoomID
	} else {
		resultRequestMessage.DestinationID = msg.OwnerID
	}

	session, err := ManagerEnv.DataCenterManager.UpdateOrCreateSession(resultRequestMessage)
	if err != nil {
		fmt.Println("Send UserMessage Error:", err)
	}

	msg.SessionID = fmt.Sprintf("%+v", session.ID)
	user, err := ManagerEnv.UserManager.GetUser(msg.OwnerID, "id")
	if err == nil {
		msg.OwnerName = user.Name
	}

	if err := ManagerEnv.DataCenterManager.Save(msg, *session); err != nil {
		return nil
	}

	// if ws, ok := m.Connections.Load(destinationID); ok {
	// 	// 3、 应该不需要将消息通过ws发送给客户端， 应该将消息包装为一个event发送给客户端
	// 	if conn, ok := ws.(*websocket.Conn); ok {
	// 		return conn.WriteJSON(struct {
	// 			Topic   string
	// 			Content interface{}
	// 		}{"message", msg})
	// 	}
	// }

	event.Pub(&MessageEvt{
		Content: msg,
	})

	// 如果用户离线，将message保存到离线数据库， redis列表的Key userID 作为唯一主键(list)
	// 当用户再次上线的时候，执行
	// 3、 应该不需要将消息通过ws发送给客户端， 应该将消息包装为一个event发送给客户端
	if err := ManagerEnv.DataCenterManager.Redis.RPush(destinationID, msg).Err(); err != nil {
		return err
	}

	return nil
}

func (m *WebsocketManager) SendRoomMessage(msg message.RequestMessage) error {
	room, err := ManagerEnv.GroupManager.GetGroup(msg.Scope.RoomID)
	if err != nil {
		return err
	}

	for _, v := range room.Childes {
		if v != msg.OwnerID {
			m.SendUserMessage(msg, v)
		}
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
