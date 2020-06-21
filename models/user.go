package models

import (
	"errors"
	"webchat/common"
	"webchat/database"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UserManager struct{}

type User struct {
	gorm.Model
	Name     string `form:"name"`
	PassWord string `form:"password"`
	Email    string
	Friends  []*User
	Room     []*Room
	Listener *Listener
}

func (user *User) BeforeSave(scope *gorm.Scope) error {
	// scope.SetColumn("ID", uuid.NewV4())
	// fmt.Println("scope beforeSave", uuid.NewV4())
	// scope.DB().Model(user).Update(user.ID, uuid.NewV4())
	scope.SetColumn("email", "827301519@qq.com")
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

func (userManager *UserManager) GetUser(ctx *gin.Context, userID string) {
	user, err := userManager.getUser(userID)
	if err != nil {
		common.Http404Response(ctx, user)
		return
	}
	common.HttpSuccessResponse(ctx, user)
}

func (*UserManager) getUser(userID string) (*User, error) {
	var users []*User
	if err := database.DB.Where("id = ?", userID).Or("name = ?", userID).Find(&users).Error; err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return users[0], nil
	}
	return nil, errors.New("record not found")
}