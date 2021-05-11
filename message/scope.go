//  用于判定发送的消息的作用域。比如是系统公告，群聊消息，个人私聊消息等
package message

type Scope struct {
	Stype         int
	RoomID        string `form:"room_id"`
	OwnerID       string `form:"owner_id"`
	DestinationID string `form:"destination_id"`
}

func NewScope(stype int, roomID, ownerID, destinationID string) *Scope {
	return &Scope{
		Stype:         stype,
		RoomID:        roomID,
		OwnerID:       ownerID,
		DestinationID: destinationID,
	}
}

func (s *Scope) GetSession() (string, error) {
	// var session models.DataCenterManager

	// if err := models.ManagerEnv.DB.Find("owner_id = ? AND destination_id = ?", s.OwnerID, s.DestinationID).Find(&session).Error; err != nil {
	// 	return "", err
	// }

	// return session.ID, nil
	return "", nil
}
