package message

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"time"
)

type MessageMedia struct {
	gorm.Model
	Scope

	StartTime time.Time
	Duration time.Duration
	Status string
}

func (m *MessageMedia) DumpMessage() string {
	m.Duration = time.Until(m.StartTime)
	buf, _ := json.Marshal(m)
	return string(buf)

}
