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

func JoinRoom(ctx *gin.Context) {

}

func LeaveRoom(ctx *gin.Context) {

}
