package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/goravel/framework/database/orm"
)

// JSONMap represents a JSON object that can be stored and retrieved from the database
type JSONMap map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONMap)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, &j)
}

type Activity struct {
	orm.Model
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Type        string       `json:"type"`
	Metadata    JSONMap      `json:"metadata" gorm:"type:json"`
	Status      string       `json:"status"`
	StartedAt   *time.Time   `json:"started_at"`
	CompletedAt *time.Time   `json:"completed_at"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"-" gorm:"index"`
}

// TableName specifies the table name for the Activity model
func (a *Activity) TableName() string {
	return "activities"
}
