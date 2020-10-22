package models

import (
	"strconv"
)

const (
	Login_Event   = "Login"
	Logout_Event  = "Logout"
	Message_Event = "Message"
)

type Event struct {
	Action                 string
	SourceID               int
	DestinationID          int
	DisplaySourceName      string
	DisplayDestinationName string
	Content                string
}

func NewEvent() *Event {
	return &Event{}
}

func (evt *Event) setDefault(action string, sourceID, destinationID int) {
	evt.Action = action
	evt.SourceID = sourceID
	evt.DestinationID = destinationID

	requestUser, _ := ManageEnv.UserManager.GetUser(strconv.Itoa(evt.SourceID))
	responseUser, _ := ManageEnv.UserManager.GetUser(strconv.Itoa(evt.DestinationID))

	evt.DisplaySourceName = requestUser.Name
	evt.DisplayDestinationName = responseUser.Name
}
