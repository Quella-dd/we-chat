package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"we-chat/models"
)

func ListSessions(c *gin.Context) {
	id := c.GetString("userID")
	sessions, err := models.ManageEnv.SessionManager.ListSessions(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
	})
}

func DeleteSession(c *gin.Context) {
	id := c.Param("id")
	if err := models.ManageEnv.SessionManager.DeleteSession(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}
