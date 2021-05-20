package models

import (
	"github.com/jinzhu/gorm"
)

type CommonManager struct{}

type Comment struct {
	gorm.Model
	OwnerID  string
	MomentID string
	Content  string
}

func NewCommonManager() *CommonManager {
	return &CommonManager{}
}

func (m *CommonManager) Create(comment Comment) error {
	return ManagerEnv.DB.Create(&comment).Error
}

func (m *CommonManager) Delete(id string) error {
	return ManagerEnv.DB.Delete(&Comment{}, id).Error
}
