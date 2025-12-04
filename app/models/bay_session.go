package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type BaySession struct {
	orm.Model
	VisitID   string    `json:"visit_id"`
	StartTime time.Time `json:"start_time"`
	Duration  int       `json:"duration"` // Duration in seconds
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for the BaySession model
func (b *BaySession) TableName() string {
	return "bay_sessions"
}
