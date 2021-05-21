package api

import (
	"net/http"
	"we-chat/models"

	"github.com/gin-gonic/gin"
)

func ListMoments(c *gin.Context) {
	userID := c.GetString("userID")
	moments, err := GE.MomentManager.List(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error)
		return
	}
	c.JSON(http.StatusOK, moments)
}

func GetMoment(c *gin.Context) {
	id := c.Param("id")
	moment, err := GE.MomentManager.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error)
		return
	}
	c.JSON(http.StatusOK, moment)
}

func CreateMoment(c *gin.Context) {
	var moment models.Moment

	err := c.ShouldBind(&moment)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = GE.MomentManager.Create(moment)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

func DeleteMoment(c *gin.Context) {
	id := c.Param("id")
	err := GE.MomentManager.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

func UpdateMoment(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")

	var spec models.UpdateSpec
	err := c.ShouldBind(&spec)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = GE.MomentManager.Update(id, userID, spec)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}
