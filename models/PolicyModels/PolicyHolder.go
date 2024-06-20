package models

import (
	"time"

	"gorm.io/gorm"
)

type PolicyHolder struct {
	gorm.Model   `swaggerignore:"true"`
	ID           uint `gorm:"primaryKey;uniqueIndex;"`
	PolicyNumber int  `gorm:"foreignKey:Policies(PolicyNumber)"`
	Person
	EffectiveDate         time.Time
	FieldChangeIndicators ChangeIndicators `gorm:"-"`
	RecordChangeIndicator string           `gorm:"-"`
}

func (PolicyHolder) TableName() string {
	return "policyprocessor_policy_holders"
}

func PullPolicyHolder(db *gorm.DB, policyNum int, effectiveDate time.Time) (*PolicyHolder, error) {
	var policyHolder PolicyHolder

	if err := db.Where("policy_number = ? AND effective_date <= ?", policyNum, effectiveDate).
		Where("deleted_at IS NULL OR deleted_at > ?", effectiveDate).
		Order("effective_date DESC, updated_at DESC").First(&policyHolder).Error; err != nil {
		return nil, err
	}

	return &policyHolder, nil
}

func PushPolicyHolder(db *gorm.DB, conditions interface{}, update interface{}) error {
	var policyholder PolicyHolder
	result := db.Where(conditions).Assign(update).FirstOrCreate(&policyholder)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
