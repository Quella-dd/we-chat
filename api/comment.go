package api

import (
	"net/http"
	"we-chat/models"

	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	var comment models.Comment
	err := c.ShouldBind(&comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = models.ManagerEnv.CommonManager.Create(comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

func DeleteComment(c *gin.Context) {
	id := c.Param("id")
	err := models.ManagerEnv.CommonManager.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}
