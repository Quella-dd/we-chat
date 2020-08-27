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
