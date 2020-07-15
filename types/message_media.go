package types

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"time"
)

type MessageMedia struct {
	ID string
	gorm.Model

	SourceID int `json: "sourceID"`
	DestinationID int `json: "destinationID"`

	StartTime time.Time
	EndTime time.Time

	Status string
}

type MessageMediaInfo struct {
	Status string
	Duration time.Duration
}

func (m *MessageMedia) DumpMessage() string {
	buf, _ := json.Marshal(struct{
		Status string
		SubTime time.Duration
	}{m.Status, m.EndTime.Sub(m.StartTime)})
	return string(buf)
}