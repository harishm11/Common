package models

import (
	"time"

	"gorm.io/gorm"
)

type Driver struct {
	gorm.Model                 `swaggerignore:"true"`
	DriverID                   int    `json:"DriverID" gorm:"primaryKey;uniqueIndex;"`
	PolicyNumber               int    `json:"PolicyNumber" gorm:"foreignKey:PolicyNumber"`
	FirstName                  string `json:"FirstName"`
	LastName                   string `json:"LastName"`
	LicenseNumber              string `json:"LicenseNumber"`
	LicenseState               string `json:"LicenseState"`
	Age                        int    `json:"Age"`
	DrivingExperience          int    `json:"DrivingExperience"`
	Gender                     string `json:"Gender"`
	MaritalStatus              string `json:"MaritalStatus"`
	DrivingCourse              string `json:"DrivingCourse"`
	MonthsSinceCourseCompleted int    `json:"MonthsSinceCourseCompleted"`
	GoodStudent                string `json:"GoodStudent"`
	StudentAway                string `json:"StudentAway"`
	EffectiveDate              time.Time
	FieldChangeIndicators      ChangeIndicators `gorm:"-"`
	RecordChangeIndicator      string           `gorm:"-"`
}

func (Driver) TableName() string {
	return "policyprocessor_drivers"
}

func PullDrivers(db *gorm.DB, policyNum int, effectiveDate time.Time) (*[]Driver, error) {
	var drivers []Driver

	// Subquery to get the latest unique records for each driver associated with the policy number
	subquery := "(SELECT DISTINCT ON (driver_id) * FROM policyprocessor_drivers WHERE policy_number = ? AND effective_date <= ? ORDER BY driver_id, effective_date DESC, updated_at DESC) AS latest_drivers"

	// Main query to filter only the active drivers
	if err := db.Raw("SELECT * FROM "+subquery+" WHERE deleted_at IS NULL OR deleted_at > ?", policyNum, effectiveDate, effectiveDate).
		Scan(&drivers).Error; err != nil {
		return nil, err
	}

	return &drivers, nil
}

func (d *Driver) GetFieldValue(field string) interface{} {
	switch field {
	case "ID":
		return d.ID
	case "DriverID":
		return d.DriverID
	case "PolicyNumber":
		return d.PolicyNumber
	case "FirstName":
		return d.FirstName
	case "LastName":
		return d.LastName
	case "LicenseNumber":
		return d.LicenseNumber
	case "LicenseState":
		return d.LicenseState
	case "Age":
		return d.Age
	case "DrivingExperience":
		return d.DrivingExperience
	case "Gender":
		return d.Gender
	case "MaritalStatus":
		return d.MaritalStatus
	case "DrivingCourse":
		return d.DrivingCourse
	case "MonthsSinceCourseCompleted":
		return d.MonthsSinceCourseCompleted
	case "GoodStudent":
		return d.GoodStudent
	case "StudentAway":
		return d.StudentAway
	case "EffectiveDate":
		return d.EffectiveDate
	default:
		return nil // Return nil for unknown fields
	}
}

func (bundle *Bundle) DriverByID(driverID int) *Driver {
	for _, driver := range bundle.Drivers {
		if driver.DriverID == driverID {
			return &driver
		}
	}
	return nil
}
