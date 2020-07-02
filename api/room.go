package api

import (
	"webchat/common"
	"webchat/models"

	"github.com/gin-gonic/gin"
)

func CreateRoom(ctx *gin.Context) {
	userID := common.GetHeader(ctx)

	var r models.Room
	if err := ctx.ShouldBind(&r); err != nil {
		ctx.JSON(400, "bad request")
		return
	}
	if err := models.ManageEnv.RoomManager.CreateRoom(userID, &r); err != nil {
		ctx.JSON(500, err)
		return
	}
	common.HttpSuccessResponse(ctx, r)
}

func ListRooms(ctx *gin.Context) {
	userID := common.GetHeader(ctx)

	if rooms, err := models.ManageEnv.RoomManager.ListRooms(userID); err != nil {
		ctx.JSON(500, err)
	} else {
		common.HttpSuccessResponse(ctx, rooms)
	}
}

func UpdateRoom(ctx *gin.Context) {
	if err := models.ManageEnv.RoomManager.UpdateRoom(ctx); err != nil {
		ctx.JSON(500, err)
	}
}

func GetRoom(ctx *gin.Context) {
	roomID := ctx.Param("id")
	if room, err := models.ManageEnv.RoomManager.GetRoom(roomID); err != nil {
		ctx.JSON(500, err)
	} else {
		common.HttpSuccessResponse(ctx, room)
	}
}

func DeleteRoom(ctx *gin.Context) {
	roomID := ctx.Param("name")
	if err := models.ManageEnv.RoomManager.DeleteRoom(roomID); err != nil {
		ctx.JSON(500, err)
	}
}
