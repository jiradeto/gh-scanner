package migrations

import (
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/jiradeto/gh-scanner/app/infrastructure/interfaces/connectors"
	"github.com/pkg/errors"
)

func getMigrationEngine() *gormigrate.Gormigrate {
	db := connectors.ConnectPostgresDB()
	migrateOpt := gormigrate.DefaultOptions
	migrateOpt.TableName = "migrations"
	migrateOpt.ValidateUnknownMigrations = true
	return gormigrate.New(db, migrateOpt, []*gormigrate.Migration{
		migrate20221215112000(),
		migrate20221215112500(),
	})
}

// Migrate is a migrate database function
func Migrate() error {
	if err := getMigrationEngine().Migrate(); err != nil {
		return errors.Wrap(err, "failed to migrate database")
	}
	log.Printf("Migration did run successfully")
	return nil
}

// RollbackLast is a rollback last migration of database function
func RollbackLast() error {
	if err := getMigrationEngine().RollbackLast(); err != nil {
		return errors.Wrap(err, "failed to rollback last migration")
	}
	log.Printf("RollbackLast did run successfully")
	return nil
}
