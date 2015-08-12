package fakes

import (
	"database/sql"

	"github.com/cloudfoundry-incubator/notifications/db"
)

type DatabaseMigrator struct {
	MigrateCall struct {
		Called   bool
		Receives struct {
			DB             *sql.DB
			MigrationsPath string
		}
	}
	SeedCall struct {
		Called   bool
		Receives struct {
			Database            db.DatabaseInterface
			DefaultTemplatePath string
		}
	}
}

func NewDatabaseMigrator() *DatabaseMigrator {
	return &DatabaseMigrator{}
}

func (d *DatabaseMigrator) Migrate(db *sql.DB, migrationsPath string) {
	d.MigrateCall.Called = true
	d.MigrateCall.Receives.DB = db
	d.MigrateCall.Receives.MigrationsPath = migrationsPath
}

func (d *DatabaseMigrator) Seed(database db.DatabaseInterface, defaultTemplatePath string) {
	d.SeedCall.Called = true
	d.SeedCall.Receives.Database = database
	d.SeedCall.Receives.DefaultTemplatePath = defaultTemplatePath
}