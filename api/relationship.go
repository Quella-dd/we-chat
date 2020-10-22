package api

import (
	"net/http"
	"we-chat/models"

	"github.com/gin-gonic/gin"
)

func GetFriends(ctx *gin.Context) {
	id := ctx.GetHeader(HeaderKey)
	if id == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	users, err := models.ManageEnv.RelationShipManager.GetFriends(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func AddFriend(ctx *gin.Context) {
	requestID := ctx.GetHeader(HeaderKey)
	destinationID := ctx.Param("id")

	if requestID == "" || destinationID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	err := models.ManageEnv.RelationShipManager.AddFriend(requestID, destinationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, nil)
}

func DeleteFriend(ctx *gin.Context) {
	requestID := ctx.GetHeader(HeaderKey)
	destinationID := ctx.Param("id")

	if requestID == "" || destinationID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	err := models.ManageEnv.RelationShipManager.DeleteFriend(requestID, destinationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, nil)
}

func ListFriendRequests(ctx *gin.Context) {
	id := ctx.GetHeader(HeaderKey)
	if id == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	requests, err := models.ManageEnv.RelationShipManager.ListAddFriendRequest(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, requests)
}

func AckFriendRequest(ctx *gin.Context) {
	requestID := ctx.GetHeader(HeaderKey)
	destinationID := ctx.Param("id")

	if requestID == "" || destinationID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	models.ManageEnv.RelationShipManager.AckRequest(requestID, destinationID)
}
