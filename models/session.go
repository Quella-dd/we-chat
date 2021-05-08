package models

import (
	"fmt"
	"time"
	Message "we-chat/message"

	"github.com/jinzhu/gorm"
)

type SessionManager struct{}

func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

type Session struct {
	gorm.Model
	Owner         string
	Src           string
	Destination   string
	LatestTime    time.Time
	LatestContent string
}

// TODO: user's icon, display
type SessionInfo struct {
	Session
	DisplayName string // TODO: displayName, 每个用户存储一份自己的数据
}

// sort with latestTime
func (s *SessionManager) ListSessions(id string) ([]SessionInfo, error) {
	var sessionInfos []SessionInfo
	var sessions []Session

	if err := ManagerEnv.DB.Where("owner = ?", id).Find(&sessions).Error; err != nil {
		return nil, nil
	}

	for _, session := range sessions {
		user, err := ManagerEnv.UserManager.GetUser(session.Destination, "id")
		if err != nil {
			fmt.Printf("user %s not found", user.Name)
		}
		sessionInfos = append(sessionInfos, SessionInfo{
			Session:     session,
			DisplayName: user.Name,
		})
	}
	return sessionInfos, nil
}

func (s *SessionManager) CreateSession(session *Session) (*Session, error) {
	if err := ManagerEnv.DB.Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionManager) GetSession(id string) ([]Message.RequestMessage, error) {
	var messages []Message.RequestMessage
	if err := ManagerEnv.DB.Where("session_id = ?", id).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// delete Session and clean the message
func (s *SessionManager) DeleteSession(id string) error {
	tx := ManagerEnv.DB.Begin()
	if err := tx.Where("id = ?", id).Delete(&Session{}).Error; err != nil {
		tx.Callback()
		return err
	}

	if err := tx.Where("session_id = ?", id).Delete(&Message.RequestMessage{}).Error; err != nil {
		tx.Callback()
		return err
	}
	tx.Commit()
	return nil
}
