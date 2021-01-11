package api

import (
	"net/http"
	"we-chat/models"

	"github.com/gin-gonic/gin"
)

// API Group
func ListGroups(c *gin.Context) {
	userID := c.GetString("userID")
	groups, err := models.ManagerEnv.GroupManager.ListGroups(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, groups)
}

func CreateGroup(c *gin.Context) {
	userID := c.GetString("userID")

	var group *models.Group
	if err := c.ShouldBind(group); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	if err := models.ManagerEnv.GroupManager.CreateGroup(userID, group); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

// update group.Name or group.Description
func UpdateGroup(c *gin.Context) {
	var group *models.Group
	if err := c.ShouldBind(group); err != nil {
		c.JSON(http.StatusBadRequest, nil)
	}
	err := models.ManagerEnv.GroupManager.UpdateGroup(group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func GetGroup(c *gin.Context) {
	id := c.Param("id")
	group, err := models.ManagerEnv.GroupManager.GetGroup(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, group)
}

func DeleteGroup(ctx *gin.Context) {
	id := ctx.Param("id")
	err := models.ManagerEnv.GroupManager.DeleteGroup(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
}

// action in group, JoinGroup and LeaveGroup
func JoinGroup(c *gin.Context) {
	id := c.GetString("userID")

	groupID := c.Param("id")
	userID := c.Param("name")

	if err := models.ManagerEnv.UserManager.JoinGroup(id, groupID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
}

func LeaveGroup(c *gin.Context) {
	id := c.GetString("userID")

	groupID := c.Param("id")
	userID := c.Param("name")

	if err := models.ManagerEnv.UserManager.LeaveGroup(id, groupID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
}