package api

import (
	"net/http"
	"we-chat/message"
	"we-chat/models"

	"github.com/gin-gonic/gin"
)

func GetMessages(c *gin.Context) {
	id := c.Param("id")
	if messages, err := models.ManagerEnv.DataCenterManager.GetMessages(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, messages)
	}
}

func HandlerMessage(c *gin.Context) {
	var requestMessage message.RequestMessage

	if err := c.ShouldBind(&requestMessage); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if err := models.ManagerEnv.DataCenterManager.HandlerMessage(requestMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func HandlerEvent(c *gin.Context) {
	id := c.GetString("userID")
	if err := models.ManagerEnv.WebsocketManager.InitWs(c, id); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, nil)
}
