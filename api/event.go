package api

import (
	"net/http"
	"we-chat/event"
	"we-chat/message"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func GetMessages(c *gin.Context) {
	id := c.Param("id")

	messagesInfos, err := GE.DataCenterManager.GetMessages(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, messagesInfos)
}

func HandlerMessage(c *gin.Context) {
	var requestMessage message.RequestMessage

	if err := c.ShouldBind(&requestMessage); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := GE.DataCenterManager.HandlerMessage(requestMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func HandlerEvent(c *gin.Context) {
	id := c.GetString("userID")

	var (
		upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
	)

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	eventCh := event.Sub(id)

	for {
		select {
		case evt := <-eventCh:
			ws.WriteJSON(evt)
		}
	}
}
