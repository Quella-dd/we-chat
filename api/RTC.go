package api

import (
	"github.com/gin-gonic/gin"
	"webchat/common"
	"webchat/models"
)

func sendRTCRequest(c *gin.Context) {
	if err := models.SendRTCRequest(c); err != nil {
		common.HttpServerError(c, err)
	}
}

func handlerRTCRequest(c *gin.Context) {
	eventID := c.Param("id")
	status := c.Param("status")

	if eventID == ""  || status == ""{
		common.HttpBadRequest(c)
		return
	}
	if err := models.HandlerRTCRequest(c, eventID, status); err != nil {
		common.HttpServerError(c, err)
	}
}
func hangDownRTCRequest(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		common.HttpBadRequest(c)
		return
	}
	if err := models.HangDownRTCRequest(c, eventID); err != nil {
		common.HttpServerError(c, err)
	}
}