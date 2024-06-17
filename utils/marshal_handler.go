package utils

import (
	"time"

	"github.com/harishm11/Common/logger"
	models "github.com/harishm11/Common/models/PolicyModels"
)

func MarshallBundle(bundle *models.Bundle, effectiveDate time.Time) (map[string]interface{}, error) {
	// Create a map to store the marshalled data
	data := make(map[string]interface{})
	gen := NewIDGenerator()
	//  Policy data
	data["LobForm"] = map[string]interface{}{
		"PolicyNumber":      bundle.Policy.PolicyNumber,
		"AccountNumber":     bundle.Policy.AccountNumber,
		"Lob":               bundle.Policy.LineofBusiness,
		"StartDate":         bundle.Policy.TermStartDate,
		"EffectiveDate":     effectiveDate,
		"TransactionType":   "",
		"TransactionNumber": int(gen.GenerateID()),
		"Quoted":            false,
	}

	//  CurrentCarrier data
	data["CurrentCarrierForm"] = map[string]interface{}{
		"PolicyNumber":      bundle.CurrentCarrier.PolicyNumber,
		"CurrentCarrier":    bundle.CurrentCarrier.CurrentCarrier,
		"CoverageStartDate": bundle.CurrentCarrier.CoverageStartDate,
		"CoverageEndDate":   bundle.CurrentCarrier.CoverageEndDate,
	}

	//  PolicyHolder data
	data["PolicyHolderForm"] = map[string]interface{}{
		"PolicyNumber": bundle.PolicyHolder.PolicyNumber,
		"FirstName":    bundle.PolicyHolder.FirstName,
		"MiddleName":   bundle.PolicyHolder.MiddleName,
		"LastName":     bundle.PolicyHolder.LastName,
		"AddressLine1": bundle.PolicyAddress.AddressLine1,
		"AddressLine2": bundle.PolicyAddress.AddressLine2,
		"City":         bundle.PolicyAddress.City,
		"State":        bundle.PolicyAddress.State,
		"ZipCode":      bundle.PolicyAddress.ZipCode,
	}

	//  Drivers data
	var driversData []interface{}
	for _, driver := range bundle.Drivers {
		driverData := map[string]interface{}{
			"DriverID":                   driver.DriverID,
			"FirstName":                  driver.FirstName,
			"LastName":                   driver.LastName,
			"LicenseNumber":              driver.LicenseNumber,
			"LicenseState":               driver.LicenseState,
			"Age":                        driver.Age,
			"DrivingExperience":          driver.DrivingExperience,
			"Gender":                     driver.Gender,
			"MaritalStatus":              driver.MaritalStatus,
			"DrivingCourse":              driver.DrivingCourse,
			"MonthsSinceCourseCompleted": driver.MonthsSinceCourseCompleted,
			"GoodStudent":                driver.GoodStudent,
			"StudentAway":                driver.StudentAway,
		}
		driversData = append(driversData, driverData)
	}
	data["DriverForm"] = driversData

	//  Vehicles data
	var vehiclesData []interface{}
	for _, vehicle := range bundle.Vehicles {
		vehicleData := map[string]interface{}{
			"VehicleID":       vehicle.VehicleID,
			"Year":            vehicle.Year,
			"Make":            vehicle.Make,
			"Model":           vehicle.ModelCd,
			"PrimaryUse":      vehicle.PrimaryUse,
			"VehicleType":     vehicle.VehicleType,
			"PrimaryOperator": vehicle.PrimaryOperator,
			"LoanORLease":     vehicle.LoanORLease,
			"Rideshare":       vehicle.Rideshare,
		}

		//  Coverages data
		var coveragesData []interface{}
		for _, coverage := range vehicle.Coverages {
			coverageData := map[string]interface{}{
				"CoverageCode":   coverage.CoverageCode,
				"CoverageOption": coverage.CoverageOption,
			}
			coveragesData = append(coveragesData, coverageData)
		}
		vehicleData["Coverages"] = coveragesData

		vehiclesData = append(vehiclesData, vehicleData)
	}
	data["VehicleForm"] = vehiclesData
	logger.GetLogger().Info("Data", data)
	return data, nil
}
