package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
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
	// Listener *Listener
}

func (user *User) BeforeSave(scope *gorm.Scope) error {
	// scope.SetColumn("ID", uuid.NewV4())
	// fmt.Println("scope beforeSave", uuid.NewV4())
	// scope.DB().Model(user).Update(user.ID, uuid.NewV4())
	scope.SetColumn("email", "827301519@qq.com")
	return nil
}

func (user User) MarshalBinary() ([]byte, error) {
	return json.Marshal(user)
}

func (user *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, user)
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

	fmt.Printf("ready to set Redis: %d, %+v\n", user.ID, user)
	if err := ManageEnv.DataCenterManager.Redis.Set(strconv.Itoa(int(user.ID)), user, time.Hour).Err(); err != nil {
		fmt.Println(err)
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

func (UserManager *UserManager) listUsers() ([]User, error){
	var users []User
	if err := database.DB.Find(&users).Error;err!= nil {
		return nil, err
	}
	return users, nil
}

//func (userManager *UserManager) JoinRoom(userID, roomID string) error {
//	user, err := userManager.getUser(userID)
//	if err != nil {
//		return err
//	}
//
//	room, err := ManageEnv.RoomManager.GetRoom(roomID)
//	if err != nil {
//		return err
//	}
//
//	if room.Childrens != nil {
//		room.Childrens = append(room.Childrens, *user)
//	} else {
//		room.Childrens = []User{*user}
//	}
//	if err := database.DB.Model(&room).Update("childrens", room.Childrens).Error; err != nil {
//		return err
//	}
//	return nil
//}
//
//func (userManager *UserManager) LeaveRoom(userID, roomID string) error {
//	_, err := userManager.getUser(userID)
//	if err != nil {
//		return err
//	}
//
//	room, err := ManageEnv.RoomManager.GetRoom(roomID)
//	if err != nil {
//		return err
//	}
//
//	if room.Childrens != nil {
//		room.filterChilds(userID)
//	} else {
//		return errors.New("this user is not in this room")
//	}
//	if err := database.DB.Model(&room).Update("childrens", room.Childrens).Error; err != nil {
//		return err
//	}
//	return nil
//}

func (userManager *UserManager) DeleteFromRoom(excuteUserID, roomID, userID string) error {
	_, err := userManager.getUser(userID)
	if err != nil {
		return err
	}

	room, err := ManageEnv.RoomManager.GetRoom(roomID)
	if err != nil {
		return err
	}

	if room.ManagerID == excuteUserID {
		room.filterChilds(userID)
		if err := database.DB.Model(&room).Update("childrens", room.Childrens).Error; err != nil {
			return err
		}
	} else {
		return errors.New("not allow to delte user from Room")
	}
	return nil
}

func (userManager *UserManager) AddUserToRoom(excuteUserID, roomID, userID string) error {
	user, err := userManager.getUser(userID)
	if err != nil {
		return err
	}

	room, err := ManageEnv.RoomManager.GetRoom(roomID)
	if err != nil {
		return err
	}

	if room.ManagerID == excuteUserID {
		if room.Childrens != nil {
			room.Childrens = append(room.Childrens, *user)
		} else {
			room.Childrens = []User{*user}
		}
		if err := database.DB.Model(&room).Update("childrens", room.Childrens).Error; err != nil {
			return err
		}
	} else {
		return errors.New("not allow to delte user from Room")
	}
	return nil
}

func (userManager *UserManager) SearchUsers(ctx *gin.Context, search string) interface{} {
	var users []User
	if err := database.DB.Where("id like?", "%"+search+"%").Or("name like ?", "%"+search+"%").Find(&users).Error; err != nil {
		fmt.Println("users sql not record")
	}
	return users
}

func (*UserManager) getUser(userID string) (*User, error) {
	var user User
	if err := database.DB.Where("id = ?", userID).Or("name = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
