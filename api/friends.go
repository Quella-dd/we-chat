package api

import (
	"webchat/common"
	"webchat/models"

	"github.com/gin-gonic/gin"
)

func GetFriends(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		common.HttpBadRequest(ctx)
		return
	}
	models.ManageEnv.RelationShipManager.GetFriends(ctx, id)
}

func AddFriend(ctx *gin.Context) {
	userID := ctx.Request.Header.Get("userID")
	id := ctx.Param("id")

	if id == "" || userID == "" {
		common.HttpBadRequest(ctx)
		return
	}
	models.ManageEnv.RelationShipManager.AddFriend(ctx, userID, id)

}

func DeleteFriend(ctx *gin.Context) {
	userID := ctx.Request.Header.Get("userID")
	id := ctx.Param("id")

	if id == "" || userID == "" {
		common.HttpBadRequest(ctx)
		return
	}
	models.ManageEnv.RelationShipManager.DeleteFriend(ctx, userID, id)
}
