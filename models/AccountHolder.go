package models

import (
	"time"

	"gorm.io/gorm"
)

type AccountHolder struct {
	gorm.Model    `swaggerignore:"true"`
	ID            uint `gorm:"primaryKey;uniqueIndex;"`
	AccountNumber int  `gorm:"foreignKey:Accounts(AccountNumber)"`
	Person
	EffectiveDate         time.Time
	FieldChangeIndicators ChangeIndicators `gorm:"-"`
	RecordChangeIndicator string           `gorm:"-"`
}

func (AccountHolder) TableName() string {
	return "policyprocessor_account_holders"
}
