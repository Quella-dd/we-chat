/*
	package scope 用户判定发送的消息的作用域。比如是系统公告，群聊消息，个人私聊消息等
*/

package message

type Scope struct {
	RoomID        int
	SourceID      int
	DestinationID int
}

func NewScope(RoomID, SourceID, DestinationID int) *Scope {
	return &Scope{
		RoomID:        RoomID,
		SourceID:      SourceID,
		DestinationID: DestinationID,
	}
}
