package models

import (
	"database/sql"
	"time"

	"github.com/goravel/framework/database/orm"
)

type BaySessionPlayer struct {
	orm.Model
	BaySessionID uint64       `json:"bay_session_id"`
	PlayerID     string       `json:"player_id"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	DeletedAt    sql.NullTime `json:"-" gorm:"index"`
}

// TableName specifies the table name for the BaySessionPlayer model
func (b *BaySessionPlayer) TableName() string {
	return "bay_session_players"
}
