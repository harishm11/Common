package migrations

import (
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// RunMigrations runs the database migrations
func RunMigrations(db *gorm.DB) {
	migrator := gormigrate.New(db, gormigrate.DefaultOptions, GetMigrations())

	if err := migrator.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Println("Migration did run successfully")
}
