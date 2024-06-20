package models

import (
	"time"

	"gorm.io/gorm"
)

type CurrentCarrierInfo struct {
	gorm.Model            `swaggerignore:"true"`
	ID                    uint `gorm:"primaryKey;uniqueIndex;"`
	PolicyNumber          int  `json:"PolicyNumber" gorm:"foreignKey:PolicyNumber"`
	EffectiveDate         time.Time
	CurrentCarrier        string           `json:"CurrentCarrierInfo" display:"Current Carrier"`
	CoverageStartDate     time.Time        `json:"CoverageStartDate" display:"Coverage Start Date"`
	CoverageEndDate       time.Time        `json:"CoverageEndDate" display:"Coverage End Date"`
	FieldChangeIndicators ChangeIndicators `gorm:"-"`
	RecordChangeIndicator string           `gorm:"-"`
}

func (CurrentCarrierInfo) TableName() string {
	return "policyprocessor_current_carrier_infos"
}

func PullCurrentCarrier(db *gorm.DB, policyNum int, effectiveDate time.Time) (*CurrentCarrierInfo, error) {
	var currentCarrier CurrentCarrierInfo

	if err := db.Where("policy_number = ? AND effective_date <= ?", policyNum, effectiveDate).
		Where("deleted_at IS NULL OR deleted_at > ?", effectiveDate).
		Order("effective_date DESC, updated_at DESC").First(&currentCarrier).Error; err != nil {
		return nil, err
	}

	return &currentCarrier, nil
}

func PushCurrentCarrier(db *gorm.DB, conditions interface{}, update interface{}) error {
	var currentCarrier CurrentCarrierInfo
	result := db.Where(conditions).Assign(update).FirstOrCreate(&currentCarrier)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (c *CurrentCarrierInfo) GetFieldValue(field string) interface{} {
	switch field {
	case "ID":
		return c.ID
	case "PolicyNumber":
		return c.PolicyNumber
	case "EffectiveDate":
		return c.EffectiveDate
	case "CurrentCarrier":
		return c.CurrentCarrier
	case "CoverageStartDate":
		return c.CoverageStartDate
	case "CoverageEndDate":
		return c.CoverageEndDate
	default:
		return nil // Return nil for unknown fields
	}
}
