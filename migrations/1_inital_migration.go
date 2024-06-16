package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/harishm11/Common/models"
	"gorm.io/gorm"
)

func GetMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "1",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Account{}, &models.Transaction{}, &models.Policy{},
					&models.AccountHolder{}, &models.AccountAddress{},
					&models.PolicyHolder{}, &models.PolicyAddress{},
					&models.Vehicle{}, &models.Driver{},
					&models.Coverage{}, &models.CurrentCarrierInfo{})

			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		// Add more migrations here
	}
}
