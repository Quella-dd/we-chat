package common

import (
	"database/sql/driver"
	"encoding/json"
)

type RelationStruct []string

func (u RelationStruct) Value() (driver.Value, error) {
	b, err := json.Marshal(u)
	return string(b), err
}

func (u *RelationStruct) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), u)
}
