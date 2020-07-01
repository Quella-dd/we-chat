package common

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const (
	USERMESSAGE = iota
	ROOMMESSAGE
	BORDERCASTMESSAGE
)

func HttpBadRequest(ctx *gin.Context) {
	ctx.JSON(400, "bad request")
}

func HttpServerError(ctx *gin.Context, errors ...error) {
	ctx.JSON(500, errors)
}

func HttpSuccessResponse(ctx *gin.Context, values ...interface{}) {
	ctx.JSON(200, values)
}

func CheckError(ctx *gin.Context, db *gorm.DB) {
	if db.Error != nil {
		HttpServerError(ctx, db.Error)
		return
	}
}

func Http404Response(ctx *gin.Context, value interface{}) {
	ctx.JSON(404, "not found")
}
