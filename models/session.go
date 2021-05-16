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
	Message.Scope
	LatestTime    time.Time
	LatestContent string
}

type SessionInfo struct {
	Session
	DisplayName string // TODO: displayName, 每个用户存储一份自己的数据
}

func (s *SessionManager) ListSessions(id string) ([]SessionInfo, error) {
	var sessionInfos []SessionInfo
	var sessions []Session

	if err := ManagerEnv.DB.Where("owner_id = ?", id).Find(&sessions).Error; err != nil {
		return nil, nil
	}

	for _, session := range sessions {
		if session.OwnerID != "" && session.DestinationID != "" {
			user, err := ManagerEnv.UserManager.GetUser(session.DestinationID, "id")
			if err != nil {
				fmt.Printf("user %s not found", user.Name)
			}
			sessionInfos = append(sessionInfos, SessionInfo{
				Session:     session,
				DisplayName: user.Name,
			})
		} else {
			sessionInfos = append(sessionInfos, SessionInfo{
				Session: session,
			})
		}
	}
	return sessionInfos, nil
}

func (s *SessionManager) CreateSession(session *Session) (*Session, error) {
	var resultSession Session
	if session.RoomID != "" {
		if err := ManagerEnv.DB.Where("owner_id = ? AND destination_id = ?", session.OwnerID, session.DestinationID).Find(&resultSession).Error; err != nil {
			if err := ManagerEnv.DB.Create(session).Error; err != nil {
				return nil, err
			}
			return session, nil
		}
	}

	if err := ManagerEnv.DB.Where("owner_id = ? AND destination_id = ? AND room_id = ?", session.OwnerID, session.DestinationID, session.RoomID).Find(&resultSession).Error; err != nil {
		if err := ManagerEnv.DB.Create(session).Error; err != nil {
			return nil, err
		}
		return session, nil
	}

	return &resultSession, nil
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
