package api

import (
	"fmt"
	"net/http"
	"time"
	"we-chat/models"

	"github.com/gin-gonic/gin"
)

func ListSessions(c *gin.Context) {
	id := c.GetString("userID")

	sessions, err := GE.SessionManager.ListSessions(id)
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
	id := c.GetString("userID")

	var session models.Session

	var sessionInfo struct {
		Stype         int
		RoomID        string
		DestinationID string
	}

	if err := c.ShouldBind(&sessionInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	fmt.Printf("create Session: %+v\n", sessionInfo)

	session.OwnerID = id

	session.Stype = sessionInfo.Stype
	session.RoomID = sessionInfo.RoomID
	session.DestinationID = sessionInfo.DestinationID

	session.LatestTime = time.Now()

	if _, err := GE.SessionManager.CreateSession(&session); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func GetSession(c *gin.Context) {
	id := c.Param("id")
	if messages, err := GE.SessionManager.GetSession(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"messages": messages,
		})
	}
}

func DeleteSession(c *gin.Context) {
	id := c.Param("id")
	if err := GE.SessionManager.DeleteSession(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}
