package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250101000005CreateBaySessionPlayersTable struct{}

// Signature The unique signature for the migration.
func (r *M20250101000005CreateBaySessionPlayersTable) Signature() string {
	return "20250101000005_create_bay_session_players_table"
}

// Up Run the migrations.
func (r *M20250101000005CreateBaySessionPlayersTable) Up() error {
	return facades.Schema().Create("bay_session_players", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("bay_session_id")
		table.UnsignedBigInteger("bay_session_location_id")
		table.String("player_id")
		table.TimestampsTz()
		table.SoftDeletesTz()
		table.Foreign("bay_session_id").References("id").On("bay_sessions")
	})
}

// Down Reverse the migrations.
func (r *M20250101000005CreateBaySessionPlayersTable) Down() error {
	return facades.Schema().DropIfExists("bay_session_players")
}
