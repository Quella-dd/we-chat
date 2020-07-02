package api

import (
	"webchat/common"
	"webchat/models"

	"github.com/gin-gonic/gin"
)

func HandlerEvent(c *gin.Context) {
	userID := c.Query("id")
	models.ManageEnv.WebsocketManager.Handler(c, userID)
}

func HandlerMessage(c *gin.Context) {
	userID := common.GetHeader(c)
	models.ManageEnv.DataCenterManager.HandlerMessage(c, userID)
}

func GetMessage(c *gin.Context) {
	userID := common.GetHeader(c)
	destID := c.Param("id")
	models.ManageEnv.DataCenterManager.GetMessage(c, userID, destID)
}
