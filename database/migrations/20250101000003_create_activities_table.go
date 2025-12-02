package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250101000003CreateActivitiesTable struct{}

// Signature The unique signature for the migration.
func (r *M20250101000003CreateActivitiesTable) Signature() string {
	return "20250101000003_create_activities_table"
}

// Up Run the migrations.
func (r *M20250101000003CreateActivitiesTable) Up() error {
	return facades.Schema().Create("activities", func(table schema.Blueprint) {
		table.ID("id")
		table.String("name")
		table.Text("description").Nullable()
		table.String("type")
		table.Json("metadata").Nullable()
		table.String("status").Default("active")
		table.Timestamp("started_at").Nullable()
		table.Timestamp("completed_at").Nullable()
		table.TimestampsTz()
		table.SoftDeletesTz()
	})
}

// Down Reverse the migrations.
func (r *M20250101000003CreateActivitiesTable) Down() error {
	return facades.Schema().DropIfExists("activities")
}
