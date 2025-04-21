package lib

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/ory/fosite"
)

// FositeArguments is a wrapper for fosite.Arguments to implement sql.Scanner and driver.Valuer.
type FositeArguments fosite.Arguments

// Scan implements the sql.Scanner interface for FositeArguments.
func (f *FositeArguments) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, f)
	case string:
		return json.Unmarshal([]byte(v), f)
	default:
		return errors.New("type assertion to []byte or string failed")
	}
}

// Value implements the driver.Valuer interface for FositeArguments.
func (f FositeArguments) Value() (driver.Value, error) {
	return json.Marshal(f)
}
