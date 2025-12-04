package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250101000004CreateBaySessionsTable struct{}

// Signature The unique signature for the migration.
func (r *M20250101000004CreateBaySessionsTable) Signature() string {
	return "20250101000004_create_bay_sessions_table"
}

// Up Run the migrations.
func (r *M20250101000004CreateBaySessionsTable) Up() error {
	return facades.Schema().Create("bay_sessions", func(table schema.Blueprint) {
		table.ID("id")
		table.String("visit_id")
		table.Timestamp("start_time")
		table.Integer("duration")
		table.TimestampsTz()
	})
}

// Down Reverse the migrations.
func (r *M20250101000004CreateBaySessionsTable) Down() error {
	return facades.Schema().DropIfExists("bay_sessions")
}
