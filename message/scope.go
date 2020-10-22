//  用于判定发送的消息的作用域。比如是系统公告，群聊消息，个人私聊消息等
package message

type Scope struct {
	Stype         int
	Ctype         string
	RoomID        int
	SourceID      int
	DestinationID int
}

func NewScope(stype int, ctype string, roomID, sourceID, destinationID int) *Scope {
	return &Scope{
		Stype:         stype,
		Ctype:         ctype,
		RoomID:        roomID,
		SourceID:      sourceID,
		DestinationID: destinationID,
	}
}
