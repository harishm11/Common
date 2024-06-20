package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model            `swaggerignore:"true"`
	AccountNumber         int
	PolicyNumber          int `gorm:"foreignKey:PolicyNumber"`
	TransactionNumber     int
	TransactionType       string `gorm:"primaryKey;uniqueIndex;"`
	Lob                   string
	EffectiveDate         time.Time
	TransactionDate       time.Time
	Status                string
	JSONRequest           string
	JSONResponse          string
	FieldChangeIndicators ChangeIndicators `gorm:"-"`
	RecordChangeIndicator string           `gorm:"-"`
}

func (Transaction) TableName() string {
	return "policyprocessor_transactions"
}

func (rr *Transaction) GetFieldValue(field string) interface{} {
	switch field {
	case "ID":
		return rr.ID
	case "AccountNumber":
		return rr.AccountNumber
	case "PolicyNumber":
		return rr.PolicyNumber
	case "TransactionNumber":
		return rr.TransactionNumber
	case "TransactionType":
		return rr.TransactionType
	case "Lob":
		return rr.Lob
	case "EffectiveDate":
		return rr.EffectiveDate
	case "Status":
		return rr.Status
	case "JSONRequest":
		return rr.JSONRequest
	case "JSONResponse":
		return rr.JSONResponse
	default:
		return nil // Return nil for unknown fields
	}
}
