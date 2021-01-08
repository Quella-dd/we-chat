package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"we-chat/database"
)

type RequestManager struct {}

func NewRequestManager() *RequestManager {
	database.DB.AutoMigrate(&Request{})
	return &RequestManager{}
}

type Request struct {
	gorm.Model
	UserID string
	AddID string
	Content string
	Status bool
}

func (m *RequestManager) CreateRequest(id, addID, content string) error {
	var request Request
	request.UserID = id
	request.AddID = addID
	request.Content = content
	request.Status = false
	return database.DB.Save(&request).Error
}

func (m *RequestManager) GetRequest(id string) (*Request, error) {
	var request Request
	if err := database.DB.Where("id = ?", id).Find(&request).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

func (m *RequestManager) ListUserRequest(id string) ([]Request, error) {
	var requests []Request
	if err := database.DB.Where("add_id = ?", id).Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

func (m *RequestManager) DeleteRequest(id string) error {
	return database.DB.Where("id = ?", id).Delete(Request{}).Error
}

func (m *RequestManager) AckRequest(id, userID string) error {
	var request Request
	if err := database.DB.Where("id = ?", id).First(&request).Error; err != nil {
		return err
	}

	if userID != request.AddID {
		return fmt.Errorf("%s not permission to resolve %s's request", userID, request.AddID)
	}

	return ManageEnv.UserManager.AckRequet(request.UserID, request.AddID)
}