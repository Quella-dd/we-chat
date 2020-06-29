package api

import (
	"webchat/models"

	"github.com/gin-gonic/gin"
)

func HandlerEvent(ctx *gin.Context) {
	userID := ctx.Query("id")
	models.ManageEnv.WebsocketManager.Handler(ctx, userID)
}

func HandlerMessage(ctx *gin.Context) {
	userID := ctx.Request.Header.Get("userID")
	models.ManageEnv.DataManager.HandlerMessage(ctx, userID)
}

func GetMessage(ctx *gin.Context) {
	userID := ctx.Request.Header.Get("userID")
	destID := ctx.Param("id")
	models.ManageEnv.DataManager.GetMessage(ctx, userID, destID)
}
