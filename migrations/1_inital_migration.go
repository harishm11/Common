package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	models "github.com/harishm11/Common/models/PolicyModels"
	"gorm.io/gorm"
)

func GetMigrations() []*gormigrate.Migration {
	var migrations []*gormigrate.Migration

	// Database-specific migrations

	migrations = getPolicyProcessorDBMigrations()

	return migrations
}

func getPolicyProcessorDBMigrations() []*gormigrate.Migration {
	// Example migrations for PolicyProcessorDB
	return []*gormigrate.Migration{
		{
			ID: "1",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(
					&models.Account{}, &models.Transaction{}, &models.Policy{},
					&models.AccountHolder{}, &models.AccountAddress{},
					&models.PolicyHolder{}, &models.PolicyAddress{},
					&models.Vehicle{}, &models.Driver{},
					&models.Coverage{}, &models.CurrentCarrierInfo{},
				)
			},
			Rollback: func(tx *gorm.DB) error {
				// Implement rollback logic if needed
				return nil
			},
		},
		// Add more migrations as needed
	}
}
