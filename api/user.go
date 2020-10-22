package api

import (
	"net/http"
	"we-chat/models"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	var user *models.User
	if err := ctx.ShouldBind(user); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	err := models.ManageEnv.UserManager.Login(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
}

func CreateUser(ctx *gin.Context) {
	var user *models.User
	if err := ctx.ShouldBind(user); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	err := models.ManageEnv.UserManager.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, user)
}

func ListUsers(ctx *gin.Context) {
	users, err := models.ManageEnv.UserManager.ListUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, "id shouldn't empty")
		return
	}
	user, err := models.ManageEnv.UserManager.GetUser(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func SearchUsers(ctx *gin.Context) {
	search := ctx.Param("search")
	if search == "" {
		ctx.JSON(http.StatusBadRequest, "search shouldn't empty")
		return
	}
	users, err := models.ManageEnv.UserManager.SearchUsers(search)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, users)
}
