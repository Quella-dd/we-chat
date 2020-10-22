package models

import (
	"encoding/json"
	"errors"
	"we-chat/database"

	"github.com/jinzhu/gorm"
)

type UserManager struct{}

type User struct {
	gorm.Model
	Name      string `form:"name"`
	PassWord  string `form:"password"`
	Email     string
	Validator bool
}

func (user *User) BeforeSave(scope *gorm.Scope) error {
	scope.DB().Model(user).Update(user.Validator, true)
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

func (*UserManager) Login(u *User) error {
	var user User
	err := database.DB.Where(u).First(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (manager *UserManager) CreateUser(user *User) error {
	if _, err := manager.getUserByName(user.Name); err == nil {
		return errors.New("duplicate name")
	}
	err := database.DB.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (*UserManager) ListUsers() ([]User, error) {
	var users []User
	err := database.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (userManager *UserManager) GetUser(id string) (*User, error) {
	var user *User
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (userManager *UserManager) SearchUsers(search string) ([]User, error) {
	var users []User
	err := database.DB.Where("id like ?", "%"+search+"%").Or("name like ?", "%"+search+"%").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (*UserManager) getUserByName(name string) (*User, error) {
	var user User
	if err := database.DB.Where("name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
