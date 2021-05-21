package api

import (
	"net/http"
	"we-chat/models"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var user models.User

	if err := c.ShouldBind(&user); err != nil {
		jsonResult(c, http.StatusBadRequest, err)
		return
	}

	u, token, err := GE.UserManager.Login(c, &user)
	if err != nil {
		jsonResult(c, http.StatusInternalServerError, err)
		return
	}

	c.Writer.Header().Set("token", token)
	jsonResult(c, http.StatusOK, u)
	return
}

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		jsonResult(c, http.StatusBadRequest, err)
		return
	}

	if err := GE.UserManager.Register(&user); err != nil {
		jsonResult(c, http.StatusInternalServerError, err)
		return
	}

	jsonResult(c, http.StatusOK, user)
	return
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, "id must not be empty")
		return
	}

	user, err := GE.UserManager.GetUser(id, "id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func SearchUsers(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, "name must not be empty")
		return
	}

	users, err := GE.UserManager.SearchUsers(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

// API Friends
func GetFriends(c *gin.Context) {
	id := c.GetString("userID")
	if id == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	users, err := GE.UserManager.ListFriends(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func AddFriend(c *gin.Context) {
	id := c.GetString("userID")
	addID := c.Param("id")

	var option models.AddUserOptions
	_ = c.ShouldBind(&option)

	if err := GE.UserManager.AddFriend(id, addID, option.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func DeleteFriend(c *gin.Context) {
	requestID := c.GetString("userID")
	destinationID := c.Param("id")

	if requestID == "" || destinationID == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	err := GE.UserManager.DeleteFriend(requestID, destinationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
