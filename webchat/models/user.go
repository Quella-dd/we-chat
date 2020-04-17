package models

import (
	"fmt"
	"webchat/common"
	"webchat/database"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type UserManager struct{}

type User struct {
	gorm.Model
	Name     string `form:"name"`
	PassWord string `form:"password"`
	Email    string
	Friends  []*User
	Room     []*Room
	// ws       *websocket.Conn
}

func (user *User) BeforeSave(scope *gorm.Scope) error {
	fmt.Printf("set uuid")
	scope.SetColumn("ID", uuid.NewV4())
	scope.DB().Model(user).Update(user.Model.ID, uuid.NewV4())
	return nil
}

func NewUserManager() *UserManager {
	database.DB.AutoMigrate(&User{})
	return &UserManager{}
}

func (*UserManager) Login(ctx *gin.Context) {
	var user, result User

	if err := ctx.ShouldBind(&user); err != nil {
		common.HttpBadRequest(ctx)
		return
	}

	if err := database.DB.Where(&user).First(&result).Error; err != nil {
		common.Http404Response(ctx, user)
	} else {
		common.HttpSuccessResponse(ctx, user)
	}
	return
}

func (*UserManager) CreateUser(ctx *gin.Context) {
	var u User
	if err := ctx.ShouldBind(&u); err != nil {
		common.HttpBadRequest(ctx)
		return
	}

	common.CheckError(ctx, database.DB.Create(&u))
	common.HttpSuccessResponse(ctx, u)
	return
}

func (*UserManager) ListUsers(ctx *gin.Context) {
	var users []User
	common.CheckError(ctx, database.DB.Find(&users))
	common.HttpSuccessResponse(ctx, users)
	return
}

func (*UserManager) GetUser(ctx *gin.Context, userID string) {
	var user User
	if database.DB.Where("id = ?", userID).First(&user).Error != nil {
		common.Http404Response(ctx, user)
		return
	}
	common.HttpSuccessResponse(ctx, user)
	return
}

func (*UserManager) GetUserByID(userID string) (*User, error) {
	var user User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (*UserManager) JoinRoom() {

}

func (*UserManager) LeaveRoom() {

}

func (*UserManager) AddFriend() {

}

func (*UserManager) RemovFriend() {

}

func (*UserManager) GetFriends() {

}
