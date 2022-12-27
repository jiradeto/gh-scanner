package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migrate20221215112500CreateRepositoryTables = []string{
	`CREATE TYPE "scan_result_status" AS ENUM ('queued', 'in_progress', 'success', 'failure')`,
	`CREATE TABLE "scan_results"
	(
		id             uuid	NOT NULL DEFAULT uuid_generate_v4(),
		repository_id  uuid	NOT NULL,
		status 		   scan_result_status NOT NULL,
		created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
		updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
		queued_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT current_timestamp,
		scanning_at    TIMESTAMP WITH TIME ZONE DEFAULT NULL,
		finished_at    TIMESTAMP WITH TIME ZONE DEFAULT NULL,
		findings       jsonb,
		PRIMARY KEY ("id"),
		CONSTRAINT "scan_result_repository_id_fkey" FOREIGN KEY ("repository_id") REFERENCES "repositories" ("id")
	)`,
}

var migrate20221215112500DropRepositoryTables = []string{
	`DROP TABLE "scan_results"`,
	`DROP TYPE "scan_result_status"`,
}

func migrate20221215112500() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20221215112500",
		Migrate: func(db *gorm.DB) error {
			for _, migrate := range migrate20221215112500CreateRepositoryTables {
				if err := db.Exec(migrate).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(db *gorm.DB) error {
			for _, migrate := range migrate20221215112500DropRepositoryTables {
				if err := db.Exec(migrate).Error; err != nil {
					return err
				}
			}
			return nil
		},
	}
}
