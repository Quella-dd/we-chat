package api

import (
	"encoding/json"
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

	var group models.Group
	if err := c.ShouldBind(&group); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	group.Childes = append(group.Childes, userID)

	if err := models.ManagerEnv.GroupManager.CreateGroup(userID, &group); err != nil {
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
	_, info := c.GetQuery("info")

	if info {
		groupInfo, err := models.ManagerEnv.GroupManager.GetGroupInfo(id, info)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, groupInfo)
		return
	}

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

	var selections struct {
		UserID string
	}
	if err := json.NewDecoder(c.Request.Body).Decode(&selections); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := models.ManagerEnv.UserManager.JoinGroup(id, groupID, selections.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
}

func LeaveGroup(c *gin.Context) {
	id := c.GetString("userID")
	groupID := c.Param("id")

	var selections struct {
		UserID string
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&selections); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := models.ManagerEnv.UserManager.LeaveGroup(id, groupID, selections.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
}
