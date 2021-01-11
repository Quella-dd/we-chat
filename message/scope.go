//  用于判定发送的消息的作用域。比如是系统公告，群聊消息，个人私聊消息等
package message

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

type Scope struct {
	Stype         int
	//Ctype         string
	RoomID        string `form:"room_id"`
	SourceID      string `form:"source_id"`
	DestinationID string `form:"destination_id"`
}

//func NewScope(stype int, ctype string, roomID, sourceID, destinationID int) *Scope {
func NewScope(stype int, roomID, sourceID, destinationID string) *Scope {
	return &Scope{
		Stype:         stype,
		//Ctype:         ctype,
		RoomID:        roomID,
		SourceID:      sourceID,
		DestinationID: destinationID,
	}
}

func (s *Scope) GetSession() (string, error) {
	hash := md5.New()
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	hash.Write(b)
	res := hash.Sum(nil)
	return hex.EncodeToString(res), nil
}
