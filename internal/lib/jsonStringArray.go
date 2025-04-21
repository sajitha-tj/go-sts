package lib

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONStringArray []string

// Scan implements the sql.Scanner interface for JSONStringArray.
func (j *JSONStringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

// Value implements the driver.Valuer interface for JSONStringArray.
func (j JSONStringArray) Value() (driver.Value, error) {
	return json.Marshal(j)
}
