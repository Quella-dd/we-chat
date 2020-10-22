package message

import (
	"encoding/json"
)

type MessageText struct {
	Content string
}

func (m *MessageText) DumpMessage() string {
	buf, _ := json.Marshal(m)
	return string(buf)
}

func NewTextMessage(content string) *MessageText {
	var msg MessageText
	if err := json.Unmarshal([]byte(content), &msg); err != nil {
		return &msg
	}
	return nil
}
