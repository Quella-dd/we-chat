package models

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

type SessionMessage struct {
	gorm.Model

	Count         int
	SourceID      int
	DestinationID int
	MessageBody   MessagesBody ` sql:"TYPE:json"`
}

type MessagesBody []string

func (sessionMessage *SessionMessage) getUserIdentify() string {
	return strconv.Itoa(sessionMessage.DestinationID)
}
