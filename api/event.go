package api

import (
	"net/http"
	"we-chat/models"

	"github.com/gin-gonic/gin"
)

func HandlerEvent(c *gin.Context) {
	userID := c.GetString("userID")
	models.ManageEnv.WebsocketManager.Handler(c, userID)
}

func HandlerMessage(c *gin.Context) {
	userID := c.GetString("userID")
	if err := models.ManageEnv.DataCenterManager.HandlerMessage(c, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func GetMessage(c *gin.Context) {
	userID := c.GetString("userID")
	destID := c.Param("id")

	if err := models.ManageEnv.DataCenterManager.GetMessage(c, userID, destID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}
