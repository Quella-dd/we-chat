package api

import (
	"webchat/common"
	"webchat/models"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	models.ManageEnv.UserManager.Login(ctx)
}

func CreateUser(ctx *gin.Context) {
	models.ManageEnv.UserManager.CreateUser(ctx)
}

func ListUsers(ctx *gin.Context) {
	models.ManageEnv.UserManager.ListUsers(ctx)
}

func GetUser(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		common.HttpBadRequest(ctx)
		return
	}
	models.ManageEnv.UserManager.GetUser(ctx, id)
}

// action in room
func RemoveFromRoom(ctx *gin.Context) {
	excuteUserID := ctx.GetHeader("userID")

	roomID := ctx.Param("id")
	userID := ctx.Param("name")
	if err := models.ManageEnv.UserManager.DeleteFromRoom(excuteUserID, roomID, userID); err != nil {
		ctx.JSON(500, err)
	}
}

func AddUserToRoom(ctx *gin.Context) {
	excuteUserID := ctx.GetHeader("userID")

	roomID := ctx.Param("id")
	userID := ctx.Param("name")
	if err := models.ManageEnv.UserManager.AddUserToRoom(excuteUserID, roomID, userID); err != nil {
		ctx.JSON(500, err)
	}
}

// common func
func SearchUsers(ctx *gin.Context) {
	search := ctx.Param("search")
	if search == "" {
		common.HttpBadRequest(ctx)
		return
	}
	obj := models.ManageEnv.UserManager.SearchUsers(ctx, search)
	ctx.JSON(200, obj)
}
