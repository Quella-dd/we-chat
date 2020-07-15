package types

import (
	"github.com/jinzhu/gorm"
)

type CommonMessageBody struct {
	Type          int `json: "Type"`

	RoomID        int `json: "RoomID"`
	SourceID      int `json: "SourceID"`
	DestinationID int `json: "DestinationID"`
}

type SessionMessage struct {
	gorm.Model

	CommonMessageBody
	Messages []interface{} ` sql:"TYPE:json"`
}

type MessageInfo struct {
	RoomID        int `json: "RoomID"`
	SourceID      int `json: "SourceID"`
	DestinationID int `json: "DestinationID"`
	Create_At     int64
	Type 		  string
	//Content       string
	Content interface{}
}


// 主要是为了用来根据Create_at来进行排序
type MessageInfos []MessageInfo

func (m MessageInfos) Len() int {
	return len(m)
}

func (m MessageInfos) Less(i, j int) bool {
	return m[i].Create_At < m[j].Create_At
}

func (m MessageInfos) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}


// 将数据库查询的数据进行转化并且返回给前端进行使用
func (sessionMessage *SessionMessage) GetMessageInfo() *[]MessageInfo {
	var messageInfos []MessageInfo

	for _, msg := range sessionMessage.Messages {
		var messageInfo MessageInfo
		messageInfo.RoomID = sessionMessage.RoomID
		messageInfo.SourceID = sessionMessage.SourceID
		messageInfo.DestinationID = sessionMessage.DestinationID

		//messageInfo.Content = msg.Content
		//messageInfo.Create_At = msg.Create_At

		messageInfo.Content = msg
		messageInfos = append(messageInfos, messageInfo)
	}

	return &messageInfos
}