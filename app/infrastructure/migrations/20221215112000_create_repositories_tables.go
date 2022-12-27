package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migrate20221215112000CreateRepositoryTables = []string{
	`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
	`CREATE TABLE "repositories"
	(
		id             uuid                     NOT NULL DEFAULT uuid_generate_v4(),
		name 		   VARCHAR                  NOT NULL,
		url 		   VARCHAR                  NOT NULL,
		created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
		updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
		PRIMARY KEY ("id")
	)`,
}

var migrate20221215112000DropRepositoryTables = []string{
	`DROP TABLE "repositories"`,
}

func migrate20221215112000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20221215112000",
		Migrate: func(db *gorm.DB) error {
			for _, migrate := range migrate20221215112000CreateRepositoryTables {
				if err := db.Exec(migrate).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(db *gorm.DB) error {
			for _, migrate := range migrate20221215112000DropRepositoryTables {
				if err := db.Exec(migrate).Error; err != nil {
					return err
				}
			}
			return nil
		},
	}
}
