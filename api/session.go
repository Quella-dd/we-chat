package api

import (
	"net/http"
	"we-chat/models"

	"github.com/gin-gonic/gin"
)

func ListSessions(c *gin.Context) {
	id := c.GetString("userID")
	sessions, err := models.ManagerEnv.SessionManager.ListSessions(id)
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

func CreateSession(c *gin.Context) {
	var session models.Session
	if err := c.ShouldBind(&session); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if _, err := models.ManagerEnv.SessionManager.CreateSession(&session); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func GetSession(c *gin.Context) {
	id := c.Param("id")
	if messages, err := models.ManagerEnv.SessionManager.GetSession(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"messages": messages,
		})
		return
	}
}

func DeleteSession(c *gin.Context) {
	id := c.Param("id")
	if err := models.ManagerEnv.SessionManager.DeleteSession(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}
