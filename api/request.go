package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// API Requests
func ListRequests(c *gin.Context) {
	id := c.GetString("userID")
	if id == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	requests, err := GE.RequestManager.ListUserRequest(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, requests)
}

func DeleteRequest(c *gin.Context) {
	id := c.Param("id")
	if err := GE.RequestManager.DeleteRequest(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func GetRequest(c *gin.Context) {
	requestID := c.Param("id")
	if requestID == "" {
		c.JSON(http.StatusBadRequest, errors.New("bad query"))
		return
	}

	requestInfo, err := GE.RequestManager.GetRequestInfo(requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, requestInfo)
}

func AckRequest(c *gin.Context) {
	requestUserID := c.GetString("userID")
	requestID := c.Param("id")

	if requestUserID == "" || requestID == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	GE.RequestManager.AckRequest(requestID, requestUserID)
}
