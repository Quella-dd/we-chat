package api

import (
	"net/http"
	"we-chat/models"

	"github.com/gin-gonic/gin"
)


// API Requests
func ListRequests(c *gin.Context) {
	id := c.GetString("userID")
	if id == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	requests, err := models.ManageEnv.RequestManager.ListUserRequest(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, requests)
}

func DeleteRequest(c *gin.Context) {
	id := c.Param("id")
	if err := models.ManageEnv.RequestManager.DeleteRequest(id); err !=  nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	} else {
		c.JSON(http.StatusOK, nil)
	}
}

func AckRequest(c *gin.Context) {
	requestID := c.GetString("userID")
	destinationID := c.Param("id")

	if requestID == "" || destinationID == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	models.ManageEnv.UserManager.AckRequet(requestID, destinationID)
}
