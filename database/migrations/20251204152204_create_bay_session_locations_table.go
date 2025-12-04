package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251204152204CreateBaySessionLocationsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251204152204CreateBaySessionLocationsTable) Signature() string {
	return "20251204152204_create_bay_session_locations_table"
}

// Up Run the migrations.
func (r *M20251204152204CreateBaySessionLocationsTable) Up() error {
	if !facades.Schema().HasTable("bay_session_locations") {
		return facades.Schema().Create("bay_session_locations", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("bay_session_id")
			table.String("bay_number")
			table.DateTime("rental_start_dt")
			table.DateTime("rental_end_dt")
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251204152204CreateBaySessionLocationsTable) Down() error {
	return facades.Schema().DropIfExists("bay_session_locations")
}
