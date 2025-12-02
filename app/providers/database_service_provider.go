package providers

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"

	"goravel/database"
)

type DatabaseServiceProvider struct {
}

func (receiver *DatabaseServiceProvider) Register(app foundation.Application) {

}

func (receiver *DatabaseServiceProvider) Boot(app foundation.Application) {
	// Initialize database if available
	if schema := facades.Schema(); schema != nil {
		kernel := database.Kernel{}
		schema.Register(kernel.Migrations())
		if seeder := facades.Seeder(); seeder != nil {
			seeder.Register(kernel.Seeders())
		}
	}
}
