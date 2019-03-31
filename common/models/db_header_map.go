package models

import (
	"encoding/json"
	"database/sql/driver"
	"errors"
	"fmt"
)

// HeaderMap Postgresql's JSONB data type
type HeaderMap map[string]interface{}

// Value get value of HeaderMap
func (hm HeaderMap) Value() (driver.Value, error) {

	//convert this map to json string, then to bytes
	jsonString, err := json.Marshal(hm)
	return []byte(jsonString), err

}

// Scan scan value into HeaderMap
func (hm *HeaderMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	return json.Unmarshal(bytes, hm)
}