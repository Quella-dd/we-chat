package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
	"we-chat/database"
)

type SessionManager struct {}

func NewSessionManager() *SessionManager {
	database.DB.AutoMigrate(&Session{})
	return &SessionManager{}
}

type Session struct {
	gorm.Model
	Owner string
	Src string
	Destination string
	LatestTime time.Time
	LatestContent interface{}
}

// TODO: user's icon, display
type SessionInfo struct {
	Session
	DisplayName string
}

// sort with latestTime
func (s *SessionManager) ListSessions(id string) ([]SessionInfo, error) {
	var sessionInfos []SessionInfo
	var sessions []Session

	if err := database.DB.Where("owner = ?", id).Find(&sessions).Error; err != nil {
		return nil, nil
	}

	for _, session := range sessions {
		user, err := ManageEnv.UserManager.GetUser(session.Destination, "id")
		if err != nil {
			fmt.Printf("user %s not found", user.Name)
		}
		sessionInfos = append(sessionInfos, SessionInfo{
			Session: session,
			DisplayName: user.Name,
		})
	}
	return sessionInfos, nil
}

func (s *SessionManager) CreateSession() error {
	return nil
}

func (s *SessionManager) DeleteSession(id string) error {
	return database.DB.Where("id = ?", id).Delete(&Session{}).Error
}