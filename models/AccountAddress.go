package models

import (
	"time"

	"gorm.io/gorm"
)

type AccountAddress struct {
	gorm.Model    `swaggerignore:"true"`
	ID            uint `gorm:"primaryKey;uniqueIndex;"`
	AccountNumber int  `gorm:"foreignKey:Accounts(AccountNumber)"`
	Address
	EffectiveDate         time.Time
	FieldChangeIndicators ChangeIndicators `gorm:"-"`
	RecordChangeIndicator string           `gorm:"-"`
}

func (AccountAddress) TableName() string {
	return "policyprocessor_account_addresses"
}
