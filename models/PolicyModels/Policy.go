package models

import (
	"time"

	"gorm.io/gorm"
)

type Policy struct {
	gorm.Model        `swaggerignore:"true"`
	PolicyNumber      int `json:"PolicyNumber" gorm:"primaryKey;uniqueIndex;"`
	AccountNumber     int `json:"AccountNumber" gorm:"index;foreignKey:AccountNumber"`
	LineofBusiness    string
	RatingCompanyCode string
	TermStartDate     time.Time
	TermEndDate       time.Time
	EffectiveDate     time.Time
	// PolicyHolder          Person           `gorm:"foreignKey:PolicyNumber"`
	// PolicyAddress         Address          `gorm:"foreignKey:PolicyNumber"`
	FieldChangeIndicators ChangeIndicators `gorm:"-"`
	RecordChangeIndicator string           `gorm:"-"`
}

func (Policy) TableName() string {
	return "policyprocessor_policies"
}

func PullPolicy(db *gorm.DB, policyNum int, effectiveDate time.Time) (*Policy, error) {
	var policy Policy

	// Retrieve the latest active policy as of the provided effective date
	if err := db.Where("policy_number = ? AND effective_date <= ?", policyNum, effectiveDate).
		Where("deleted_at IS NULL OR deleted_at > ?", effectiveDate).
		Order("effective_date DESC, updated_at DESC").First(&policy).Error; err != nil {
		return nil, err
	}

	return &policy, nil
}

func (p *Policy) GetFieldValue(field string) interface{} {
	switch field {
	case "PolicyNumber":
		return p.PolicyNumber
	case "AccountNumber":
		return p.AccountNumber
	case "LineofBusiness":
		return p.LineofBusiness
	case "RatingCompanyCode":
		return p.RatingCompanyCode
	case "TermStartDate":
		return p.TermStartDate
	case "TermEndDate":
		return p.TermEndDate
	case "EffectiveDate":
		return p.EffectiveDate
	// Add cases for other fields if needed...
	default:
		return nil // Return nil if the field is not found
	}
}
