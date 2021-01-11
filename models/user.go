package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type UserManager struct{}

func NewUserManager() *UserManager {
	return &UserManager{}
}

type User struct {
	gorm.Model
	Name      string `form:"name"`
	PassWord  string `form:"password"`
	Email     string `form:"email"`
	Validate bool
	Relations RelationStruct `gorm:"type:json"`
}

type AddUserOptions struct {
	Content string `form:"content"`
}

type RelationStruct []string

func (u RelationStruct) Value() (driver.Value, error) {
	b, err := json.Marshal(u)
	return string(b), err
}

func (u *RelationStruct) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), u)
}

func (*UserManager) Login(u *User) (string, error) {
	var user User

	// TODO: 用户密码使用MD5进行解密并且验证
	//user.PassWord = md5Password(user.PassWord)

	err := ManagerEnv.DB.Where(u).First(&user).Error
	if err != nil {
		return "", err
	}
	token, err := GenerateToken(user.Name, strconv.Itoa(int(user.ID)), 24 * time.Hour)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (m *UserManager) Register(user *User) error {
	if _, err := m.GetUser(user.Name, "name"); err == nil {
		return errors.New("username must not be duplicate")
	}

	// TODO: 用户密码使用MD5进行加密并保存在数据库中
	//user.PassWord = md5Password(user.PassWord)
	if err := ManagerEnv.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (*UserManager) ListUsers() ([]User, error) {
	var users []User
	err := ManagerEnv.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (m *UserManager) GetUser(search, option string) (*User, error) {
	var user User
	query := fmt.Sprintf("%s = ?", option)
	if err := ManagerEnv.DB.Where(query, search).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserManager) SearchUsers(name string) ([]User, error) {
	var users []User
	query := fmt.Sprintf("%%%s%%", name)
	if err := ManagerEnv.DB.Where("name LIKE ?", query).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// ----------------------- Friends
func (m *UserManager) ListFriends(id string) ([]User, error) {
	var users []User
	user, err := m.GetUser(id, "id")
	if err != nil {
		return nil, err
	}
	if err := ManagerEnv.DB.Find(&users, []string(user.Relations)).Error; err != nil {
		return nil, fmt.Errorf("ListFriend error:", err)
	}
	return users, nil
}

func (m *UserManager) AddFriend(id, addID, content string) error {
	if content == "" {
		user, err := ManagerEnv.UserManager.GetUser(addID, "id")
		if err != nil {
			return err
		}
		content = fmt.Sprintf("'I'm %s\n", user.Name)
	}
	return ManagerEnv.RequestManager.CreateRequest(id, addID, content)
}

func (m *UserManager) DeleteFriend(id, friendID string) error {
	self, _:= m.GetUser(id, "id")

	if _,  err := m.GetUser(friendID, "id"); err != nil {
		return errors.New("user not found")
	}
	for i, v := range self.Relations {
		if v == friendID {
			self.Relations = append(self.Relations[:i], self.Relations[i+1:]...)
		}
	}
	return ManagerEnv.DB.Model(&self).Update("relations", self.Relations).Error
}

func (m *UserManager) AckRequet(id, friendID string) error {
	self, _:= m.GetUser(id, "id")

	if _,  err := m.GetUser(friendID, "id"); err != nil {
		return errors.New("user not found")
	}
	self.Relations = append(self.Relations, friendID)

	// 更新user relations and request's status, 因此要开启一个事务，保证数据一致性
	tx := ManagerEnv.DB.Begin()
	if err := tx.Model(&self).Update("relations", self.Relations).Error; err != nil {
		tx.Rollback()
		return err
	}
	request, err := ManagerEnv.RequestManager.GetRequest(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&request).Updates(Request{Status: true}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}