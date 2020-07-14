package api

import (
	"webchat/common"
	"webchat/models"

	"github.com/gin-gonic/gin"
)

func GetFriends(ctx *gin.Context) {
	userID := ctx.Request.Header.Get("userID")
	if userID == "" {
		common.HttpBadRequest(ctx)
		return
	}
	models.ManageEnv.RelationShipManager.GetFriends(ctx, userID)
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


func ListFriendRequests(ctx *gin.Context) {
	userID := ctx.Request.Header.Get("userID")
	if userID == ""  {
		common.HttpBadRequest(ctx)
		return
	}
	models.ManageEnv.RelationShipManager.ListRequest(ctx, userID)
}

func AckFriendRequest(ctx *gin.Context) {
	userID := ctx.Request.Header.Get("userID")
	id := ctx.Param("id")


	if userID == ""  || id == "" {
		common.HttpBadRequest(ctx)
		return
	}
	models.ManageEnv.RelationShipManager.AckRequest(ctx, userID, id)
}