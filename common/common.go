package common

import (
	"github.com/gin-gonic/gin"
)

const (
	HeaderKey = "userID"
)

const (
	USERMESSAGE = iota 	// 私聊时的文本消息
	RTCMESSAGE			// 私聊时的实时视频消息
	ROOMMESSAGE			// 群聊时的文本消息
	BORDERCASTMESSAGE	// 类似系统公告
)

const (
	UnConfirmd = 0
	Confirmd = 1
)

func GetHeader(c *gin.Context) string {
	return c.GetHeader(HeaderKey)
}
