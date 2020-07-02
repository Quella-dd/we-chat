package common

import (
	"github.com/gin-gonic/gin"
)

const (
	USERMESSAGE = iota
	ROOMMESSAGE
	BORDERCASTMESSAGE
)

const (
	HeaderKey = "userID"
)

func GetHeader(c *gin.Context) string {
	return c.GetHeader(HeaderKey)
}
