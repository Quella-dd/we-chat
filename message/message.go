/*
	package Message 主要用来实现各种message的定义， 并且最终序列化到是数据库中，作为用户的一条message或者存储在redis缓存中

	消息内容包含：
		1、文本消息
		2、图片消息
		3、视频请求消息 (webRTC）
//
		4、音频请求消息 (webRTC）
		5、语音消息

	MessageText
	MessageImage

	(在序列化为message时，需要延迟执行， 需要记录消息的状态，成功与否，如果成功则需记录消息的时长)
	MessageVedioOnline
	MessageAudioOnline

	MessageAudioOffline

	interface Marshal {
		Marshal()
	}
*/

package message

import (
	"github.com/jinzhu/gorm"
)

const (
	MessageSignalPrivate = "private"
	MessageSignalGroup   = "group"
	MessageSignalPublic  = "public"
)

// 用户发起的message
type RequestMessage struct {
	Singal string `json: "Singal"`
	Type   int    `json: "Type"`

	RoomID        int    `json: "RoomID"`
	SourceID      int    `json: "SourceID"`
	DestinationID int    `json: "DestinationID"`
	Body          string `json: "Body`
}

// 数据库中存的数据
type SessionMessage struct {
	gorm.Model
	RoomID        int           `json: "RoomID"`
	SourceID      int           `json: "SourceID"`
	DestinationID int           `json: "DestinationID"`
	Messages      []interface{} ` sql:"TYPE:json"`
}

// type MessageInfo struct {
// 	RoomID        int `json: "RoomID"`
// 	SourceID      int `json: "SourceID"`
// 	DestinationID int `json: "DestinationID"`
// 	Create_At     int64
// 	Type          string
// 	//Content       string
// 	Content interface{}
// }

// // 主要是为了用来根据Create_at来进行排序
// type MessageInfos []MessageInfo

// func (m MessageInfos) Len() int {
// 	return len(m)
// }

// func (m MessageInfos) Less(i, j int) bool {
// 	return m[i].Create_At < m[j].Create_At
// }

// func (m MessageInfos) Swap(i, j int) {
// 	m[i], m[j] = m[j], m[i]
// }

// // GetMessageInfo 将数据库查询的数据进行转化并且返回给前端进行使用
// func (sessionMessage *SessionMessage) GetMessageInfo() *[]MessageInfo {
// 	var messageInfos []MessageInfo

// 	for _, msg := range sessionMessage.Messages {
// 		var messageInfo MessageInfo
// 		messageInfo.RoomID = sessionMessage.RoomID
// 		messageInfo.SourceID = sessionMessage.SourceID
// 		messageInfo.DestinationID = sessionMessage.DestinationID

// 		//messageInfo.Content = msg.Content
// 		//messageInfo.Create_At = msg.Create_At

// 		messageInfo.Content = msg
// 		messageInfos = append(messageInfos, messageInfo)
// 	}

// 	return &messageInfos
// }

// // NewMessage
// func NewMessage() interface{} {
// 	return nil
// }

// /*
// 	Signal string // 私聊，群聊，系统公告
// 	Type string // 文字，视频，语音。。。

// 	SourceID string
// 	DestinationID string
// 	Content string
// */
