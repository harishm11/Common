package models

import (
	"time"

	"gorm.io/gorm"
)

type Vehicle struct {
	gorm.Model            `swaggerignore:"true"`
	ID                    uint `gorm:"primaryKey;uniqueIndex;"`
	VehicleID             int  `json:"VehicleID"`
	VIN                   string
	PolicyNumber          int        `json:"PolicyNumber" gorm:"foreignKey:PolicyNumber"`
	Year                  int        `json:"Year"`
	Make                  string     `json:"Make"`
	ModelCd               string     `json:"Model"`
	PrimaryUse            string     `json:"PrimaryUse" display:"Primary Use"`
	VehicleType           string     `json:"VehicleType" display:"Vehicle Type"`
	PrimaryOperator       string     `json:"PrimaryOperator" display:"Primary Operator"`
	LoanORLease           string     `json:"LoanORLease"`
	Rideshare             string     `json:"Rideshare"`
	Coverages             []Coverage `json:"Coverages" gorm:"foreignKey:VehicleID" `
	EffectiveDate         time.Time
	FieldChangeIndicators ChangeIndicators `gorm:"-"`
	RecordChangeIndicator string           `gorm:"-"`
}

func (Vehicle) TableName() string {
	return "policyprocessor_vehicles"
}

func PullVehicles(db *gorm.DB, policyNum int, effectiveDate time.Time) (*[]Vehicle, error) {
	var vehicles []Vehicle

	// Subquery to get the latest unique records for each vehicle as of the effective date passed
	subquery := "(SELECT DISTINCT ON (vehicle_id) * FROM policyprocessor_vehicles WHERE policy_number = ? AND effective_date <= ? ORDER BY vehicle_id, effective_date DESC, updated_at DESC) AS latest_vehicles"

	// Main query to filter only the active vehicles
	if err := db.Raw("SELECT * FROM "+subquery+" WHERE deleted_at IS NULL OR deleted_at > ?", policyNum, effectiveDate, effectiveDate).
		Scan(&vehicles).Error; err != nil {
		return nil, err
	}

	return &vehicles, nil
}

func (v *Vehicle) GetFieldValue(field string) interface{} {
	switch field {
	case "ID":
		return v.ID
	case "VehicleID":
		return v.VehicleID
	case "VIN":
		return v.VIN
	case "PolicyNumber":
		return v.PolicyNumber
	case "Year":
		return v.Year
	case "Make":
		return v.Make
	case "ModelCd":
		return v.ModelCd
	case "PrimaryUse":
		return v.PrimaryUse
	case "VehicleType":
		return v.VehicleType
	case "PrimaryOperator":
		return v.PrimaryOperator
	case "LoanORLease":
		return v.LoanORLease
	case "Rideshare":
		return v.Rideshare
	case "EffectiveDate":
		return v.EffectiveDate
	default:
		return nil
	}
}

func (bundle *Bundle) VehicleByID(vehicleID int) *Vehicle {
	for _, vehicle := range bundle.Vehicles {
		if vehicle.VehicleID == vehicleID {
			return &vehicle
		}
	}
	return nil
}

// CoverageByID returns the coverage with the specified ID from the given vehicle.
func (v *Vehicle) CoverageByID(CoverageCode string) *Coverage {
	for _, coverage := range v.Coverages {
		if coverage.CoverageCode == CoverageCode {
			return &coverage
		}
	}
	return nil
}
