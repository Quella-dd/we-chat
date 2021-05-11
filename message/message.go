package message

import (
	"github.com/jinzhu/gorm"
)

const (
	USERMESSAGE       = iota // 私聊时的文本
	RTCMESSAGE               // 私聊时的实时语言/视频
	ROOMMESSAGE              // 群聊时的文本
	BORDERCASTMESSAGE        // 系统公告
)

const (
	Message_Text   = "message_text"
	Message_Signal = "message_signal"
	Message_Media  = "message_media"
	Message_Image  = "message_image"
)

type RequestMessage struct {
	gorm.Model
	OwnerName string
	SessionID string
	Content   string

	Scope
}

const (
	MessageSignalPrivate = "private"
	MessageSignalGroup   = "group"
	MessageSignalPublic  = "public"
)

// 每种message 都有自己独特的序列化的方法，前端传递过来的结构中，将各自的字段转换为string类型的字符传递过来，
type MessageInterface interface {
	Marshal()
	Unmarshal()
}
