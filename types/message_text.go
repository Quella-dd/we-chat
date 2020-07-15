package types

import (
	"encoding/json"
)

type MessageText struct {

}

func (m *MessageText) DumpMessage() string {
	buf, _ := json.Marshal(m)
	return string(buf)
}
