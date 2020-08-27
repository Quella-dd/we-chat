package api

import (
	"net/http"
	"webchat/models"

	"github.com/gin-gonic/gin"
)

const (
	HeaderKey = "userID"
)

func ListGroup(ctx *gin.Context) {
	userID := ctx.GetHeader(HeaderKey)
	rooms, err := models.ManageEnv.GroupManager.ListGroup(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, rooms)
}

func CreateGroup(ctx *gin.Context) {
	var group *models.Group
	userID := ctx.GetHeader(HeaderKey)

	if err := ctx.ShouldBind(group); err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	if err := models.ManageEnv.GroupManager.CreateGroup(userID, group); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

// update group.Name or group.Description
func UpdateGroup(ctx *gin.Context) {
	var group *models.Group
	if err := ctx.ShouldBind(group); err != nil {
		ctx.JSON(http.StatusBadRequest,nil)
	}
	err := models.ManageEnv.GroupManager.UpdateGroup(group)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func GetGroup(ctx *gin.Context) {
	id := ctx.Param("id")
	if room, err := models.ManageEnv.GroupManager.GetGroup(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	} else {
		ctx.JSON(http.StatusOK, room)
	}
}

func DeleteGroup(ctx *gin.Context) {
	id := ctx.Param("id")
	err := models.ManageEnv.GroupManager.DeleteGroup(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
}

// action in group
func AddUserToRoom(ctx *gin.Context) {
	id := ctx.GetHeader("userID")
	groupID := ctx.Param("id")
	userID := ctx.Param("name")

	if err := models.ManageEnv.UserManager.AddUserToGroup(id, groupID, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
}

func RemoveFromRoom(ctx *gin.Context) {
	id := ctx.GetHeader("userID")
	groupID := ctx.Param("id")
	userID := ctx.Param("name")

	if err := models.ManageEnv.UserManager.DeleteFromGroup(id, groupID, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
}
