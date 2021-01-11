package api

import (
	"net/http"
	"we-chat/message"
	"we-chat/models"

	"github.com/gin-gonic/gin"
)

func HandlerMessage(c *gin.Context) {
	var requestMessage message.RequestMessage

	if err := c.ShouldBind(&requestMessage); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if err := models.ManagerEnv.DataCenterManager.HandlerMessage(c, requestMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}