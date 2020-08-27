package message

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"strconv"
)

const (
	USERMESSAGE = iota 	// 私聊时的文本消息
	RTCMESSAGE			// 私聊时的实时视频消息
	ROOMMESSAGE			// 群聊时的文本消息
	BORDERCASTMESSAGE	// 类似系统公告
)

const (
	MessageSignalPrivate = "private"
	MessageSignalGroup   = "group"
	MessageSignalPublic  = "public"
)

// 每种message 都有自己独特的序列化的方法，前端传递过来的结构中，将各自的字段转换为string类型的字符传递过来，
type DumpRedisMessage interface {
	DumpMessage() string
}

type RequestMessage struct {
	Scope
	Body          string `json: "Body`
}

type SessionMessage struct {
	gorm.Model
	Scope

	Messages      []interface{} `sql:"TYPE:json"`
}

type MessageInfo struct {
	RoomID        int `json: "RoomID"`
	SourceID      int `json: "SourceID"`
	DestinationID int `json: "DestinationID"`
	Create_At     int64
	Type          string
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
	Create_At     int64
	Type          string
	Content       string
	//Value interface{}
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
