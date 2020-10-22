package api

import (
	"net/http"
	"we-chat/models"

	"github.com/gin-gonic/gin"
)

func sendRTCRequest(c *gin.Context) {
	if err := models.SendRTCRequest(c); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
}

func handlerRTCRequest(c *gin.Context) {
	eventID := c.Param("id")
	status := c.Param("status")

	if eventID == "" || status == "" {
		c.JSON(http.StatusBadRequest, eventID)
		return
	}
	if err := models.HandlerRTCRequest(c, eventID, status); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
}

func hangDownRTCRequest(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.JSON(http.StatusBadRequest, eventID)
		return
	}
	if err := models.HangDownRTCRequest(c, eventID); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
}
