package models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model    `swaggerignore:"true"`
	AccountNumber int `json:"AccountNumber" gorm:"primaryKey;uniqueIndex;" `
	//AccountHolder  *AccountHolder  `gorm:"foreignKey:AccountNumber"` // hasOne relationship
	//AccountAddress *AccountAddress `gorm:"foreignKey:AccountNumber"` // hasOne relationship
}

func (Account) TableName() string {
	return "policyprocessor_accounts"
}
