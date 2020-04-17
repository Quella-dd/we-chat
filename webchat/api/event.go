package api

import (
	"fmt"
	"webchat/models"

	"github.com/gin-gonic/gin"
)

func HandlerEvent(ctx *gin.Context) {
	userID := ctx.Query("id")
	fmt.Println(userID)
	models.ManageEnv.WebsocketManager.Handler(ctx, userID)
}

func HandlerMessage(ctx *gin.Context) {
	models.ManageEnv.DataManager.HandlerMessage(ctx)
}
