package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	models "github.com/harishm11/Common/models/PolicyModels"
	ratingmodels "github.com/harishm11/Common/models/RatingModels"
	"gorm.io/gorm"
)

func GetMigrations(dbName string) []*gormigrate.Migration {
	var migrations []*gormigrate.Migration

	// Database-specific migrations
	switch dbName {
	case "PolicyProcessorDB":
		migrations = getPolicyProcessorDBMigrations()
	case "RateDB":
		migrations = getRateDBMigrations()
		// Add more cases as needed for different databases
	}

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

func getRateDBMigrations() []*gormigrate.Migration {
	// Example migrations for RateDB
	return []*gormigrate.Migration{
		{
			ID: "2",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&ratingmodels.Ratebooks{},
					&ratingmodels.RateFactors{},
					&ratingmodels.RateRoutines{},
					&ratingmodels.RateRoutinSteps{},
				)
			},
			Rollback: func(tx *gorm.DB) error {
				// Rollback logic specific to RateDB
				return nil
			},
		},
	}
}
