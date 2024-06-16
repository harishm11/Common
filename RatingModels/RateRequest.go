package models

import "time"

type RateRequest struct {
	AccountNumber     int
	PolicyNumber      int
	TransactionNumber int
	EffectiveDate     time.Time
	TransactionDate   time.Time
	VehDetails        []ReqVehData
	DrvDetails        []ReqDrvData
}
type ReqCvgData struct {
	ID             uint   `gorm:"primaryKey;uniqueIndex;"`
	VehicleID      uint   `gorm:"foreignKey:ID"`
	CoverageCode   string `json:"CoverageCode" display:"Coverage Code"`
	CvgSymbol      string
	CoverageOption string `json:"CoverageOption" display:"Coverage Option"`
}

type ReqVehData struct {
	ID              uint `gorm:"primaryKey;uniqueIndex;"`
	VIN             string
	PolicyNumber    int          `json:"PolicyNumber" gorm:"foreignKey:PolicyNumber"`
	Year            int          `json:"Year"`
	Make            string       `json:"Make"`
	ModelCd         string       `json:"Model"`
	PrimaryUse      string       `json:"PrimaryUse" display:"Primary Use"`
	VehicleType     string       `json:"VehicleType" display:"Vehicle Type"`
	PrimaryOperator string       `json:"PrimaryOperator" display:"Primary Operator"`
	LoanORLease     string       `json:"LoanORLease"`
	Rideshare       string       `json:"Rideshare"`
	CvgDetails      []ReqCvgData `json:"Coverages" gorm:"foreignKey:VehicleID" `
}

type ReqDrvData struct {
	ID                         uint   `gorm:"primaryKey;uniqueIndex;"`
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
}
