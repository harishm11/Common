package models

import (
	"time"

	"gorm.io/gorm"
)

type PolicyAddress struct {
	gorm.Model   `swaggerignore:"true"`
	ID           uint `gorm:"primaryKey;uniqueIndex;"`
	PolicyNumber int  `gorm:"foreignKey:Policies(PolicyNumber)"`
	Address
	EffectiveDate         time.Time
	FieldChangeIndicators ChangeIndicators `gorm:"-"`
	RecordChangeIndicator string           `gorm:"-"`
}

func (PolicyAddress) TableName() string {
	return "policyprocessor_policy_addresses"
}

func PullPolicyAddress(db *gorm.DB, policyNum int, effectiveDate time.Time) (*PolicyAddress, error) {
	var policyAddress PolicyAddress

	if err := db.Where("policy_number = ? AND effective_date <= ?", policyNum, effectiveDate).
		Where("deleted_at IS NULL OR deleted_at > ?", effectiveDate).
		Order("effective_date DESC, updated_at DESC").First(&policyAddress).Error; err != nil {
		return nil, err
	}

	return &policyAddress, nil
}
