package models

import (
	"time"

	"gorm.io/gorm"
)

type Coverage struct {
	gorm.Model            `swaggerignore:"true"`
	ID                    uint   `gorm:"primaryKey;uniqueIndex;"`
	VehicleID             uint   `gorm:"foreignKey:ID"`
	CoverageCode          string `json:"CoverageCode" display:"Coverage Code"`
	CvgSymbol             string
	CoverageOption        string `json:"CoverageOption" display:"Coverage Option"`
	CvgPremium            float32
	EffectiveDate         time.Time
	FieldChangeIndicators ChangeIndicators `gorm:"-"`
	RecordChangeIndicator string           `gorm:"-"`
}

func (Coverage) TableName() string {
	return "policyprocessor_coverages"
}

func PullCoveragesForVehicle(db *gorm.DB, vehicleID uint, effectiveDate time.Time) (*[]Coverage, error) {
	var coverages []Coverage

	// Subquery to get the latest unique records for each coverage associated with the vehicle ID
	subquery := "(SELECT DISTINCT ON (coverage_code) * FROM policyprocessor_coverages WHERE vehicle_id = ? AND effective_date <= ? ORDER BY coverage_code, effective_date DESC, updated_at DESC) AS latest_coverages"

	// Main query to filter only the active coverages
	if err := db.Raw("SELECT * FROM "+subquery+" WHERE deleted_at IS NULL OR deleted_at > ?", vehicleID, effectiveDate, effectiveDate).
		Scan(&coverages).Error; err != nil {
		return nil, err
	}

	return &coverages, nil
}

func (c *Coverage) GetFieldValue(field string) interface{} {
	switch field {
	case "ID":
		return c.ID
	case "VehicleID":
		return c.VehicleID
	case "CoverageCode":
		return c.CoverageCode
	case "CvgSymbol":
		return c.CvgSymbol
	case "CoverageOption":
		return c.CoverageOption
	case "CvgPremium":
		return c.CvgPremium
	case "EffectiveDate":
		return c.EffectiveDate
	default:
		return nil // Return nil for unknown fields
	}
}
