package models

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/jinzhu/gorm"
)

type CommonMessageBody struct {
	Type          int `json: "Type"`
	Count         int `json: "Count"`
	RoomID        int `json: "RoomID"`
	SourceID      int `json: "SourceID"`
	DestinationID int `json: "DestinationID"`
}

type RequestBody struct {
	CommonMessageBody
	Content  string `json: "Content"`
	CreateAt int64  `json: "CreateAt"`
}

type SessionMessage struct {
	gorm.Model

	CommonMessageBody
	Messages SliceMessageBody ` sql:"TYPE:json"`
}

type SliceMessageBody []MessagesBody

type MessagesBody struct {
	Create_At int64
	Content   string
}

type MessageInfo struct {
	RoomID        int `json: "RoomID"`
	SourceID      int `json: "SourceID"`
	DestinationID int `json: "DestinationID"`
	Content       string
	Create_At     int64
}

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

func (sessionMessage *SessionMessage) GetMessageInfo() *[]MessageInfo {
	var messageInfos []MessageInfo

	for _, msg := range sessionMessage.Messages {
		var messageInfo MessageInfo
		messageInfo.RoomID = sessionMessage.RoomID
		messageInfo.SourceID = sessionMessage.SourceID
		messageInfo.DestinationID = sessionMessage.DestinationID
		messageInfo.Content = msg.Content
		messageInfo.Create_At = msg.Create_At
		messageInfos = append(messageInfos, messageInfo)
	}

	return &messageInfos
}

func (c SliceMessageBody) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *SliceMessageBody) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

func (msg RequestBody) MarshalBinary() (data []byte, err error) {
	return json.Marshal(msg)
}

func (msg *RequestBody) UnmarshalBinary(data []byte) error {
	var buf []byte
	return json.Unmarshal(buf, &msg)
}

func (msg *RequestBody) getUserIdentify() string {
	return strconv.Itoa(msg.DestinationID)
}
