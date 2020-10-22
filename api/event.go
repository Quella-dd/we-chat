package api

import (
	"we-chat/models"

	"github.com/gin-gonic/gin"
)

func HandlerEvent(c *gin.Context) {
	userID := c.Query("id")
	models.ManageEnv.WebsocketManager.Handler(c, userID)
}

func HandlerMessage(c *gin.Context) {
	userID := c.GetHeader(HeaderKey)
	models.ManageEnv.DataCenterManager.HandlerMessage(c, userID)
}

func GetMessage(c *gin.Context) {
	userID := c.GetHeader(HeaderKey)
	destID := c.Param("id")

	if err := models.ManageEnv.DataCenterManager.GetMessage(c, userID, destID); err != nil {
		jsonResult(c.Writer, err)
	}
}
