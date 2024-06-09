package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/harishm11/API-Gateway/logger"
	"github.com/harishm11/API-Gateway/models"
)

func UnmarshallRequest(c *fiber.Ctx, data map[string]interface{}, bundle *models.Bundle) error {
	defer func() {
		if r := recover(); r != nil {
			// Recover from panic and send a Fiber error in the context
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}
	}()
	//  Policy data
	policyData := data["LobForm"].(map[string]interface{})
	gen := NewIDGenerator()
	policyNumber := int(policyData["PolicyNumber"].(float64))
	lob := policyData["Lob"]
	var DateValue string
	FirstQuote := policyData["TransactionType"] == "Submission" && !policyData["Quoted"].(bool)
	if FirstQuote {
		policyNumber = int(gen.GenerateID())
	}
	bundle.Policy = models.Policy{
		PolicyNumber:   policyNumber,
		AccountNumber:  int(policyData["AccountNumber"].(float64)),
		LineofBusiness: lob.(string),
	}

	DateValue = policyData["EffectiveDate"].(string)
	eff_date := ParseDate(DateValue)
	bundle.Policy.EffectiveDate = eff_date

	DateValue = policyData["StartDate"].(string)
	start_date := ParseDate(DateValue)
	if policyData["TransactionType"].(string) == "Renewal" {
		bundle.Policy.TermStartDate = eff_date
	} else {
		bundle.Policy.TermStartDate = start_date
	}

	if policyData["TransactionType"].(string) == "Cancel" {
		bundle.Policy.TermEndDate = bundle.Policy.EffectiveDate
	} else {
		bundle.Policy.TermEndDate = bundle.Policy.TermStartDate.AddDate(0, 6, 0)
	}

	//  Current Carrier data
	currentCarrierData := data["CurrentCarrierForm"].(map[string]interface{})
	DateValue = currentCarrierData["CoverageStartDate"].(string)
	cvgst_date := ParseDate(DateValue)
	DateValue = currentCarrierData["CoverageEndDate"].(string)
	cvgend_date := ParseDate(DateValue)
	bundle.CurrentCarrier = models.CurrentCarrierInfo{
		PolicyNumber:      policyNumber,
		EffectiveDate:     eff_date,
		CurrentCarrier:    currentCarrierData["CurrentCarrier"].(string),
		CoverageStartDate: cvgst_date,
		CoverageEndDate:   cvgend_date,
	}

	//PolicyHolder data
	policyHolderData := data["PolicyHolderForm"].(map[string]interface{})
	bundle.PolicyHolder = models.PolicyHolder{
		PolicyNumber:  policyNumber,
		EffectiveDate: eff_date,
		Person: models.Person{
			FirstName:  policyHolderData["FirstName"].(string),
			LastName:   policyHolderData["MiddleName"].(string),
			MiddleName: policyHolderData["LastName"].(string),
		},
	}
	bundle.PolicyAddress = models.PolicyAddress{
		PolicyNumber:  policyNumber,
		EffectiveDate: eff_date,
		Address: models.Address{
			AddressLine1: handleNullString(policyHolderData["AddressLine1"]),
			AddressLine2: handleNullString(policyHolderData["AddressLine2"]),
			City:         handleNullString(policyHolderData["City"]),
			County:       handleNullString(policyHolderData["County"]),
			State:        handleNullString(policyHolderData["State"]),
			ZipCode:      handleNullString(policyHolderData["ZipCode"]),
		},
	}

	// Drivers data
	driversData := data["DriverForm"].([]interface{})
	var drivers []models.Driver
	for _, driver := range driversData {
		driverData := driver.(map[string]interface{})

		driverNumber := int(driverData["DriverID"].(float64))

		d := models.Driver{
			EffectiveDate:     eff_date,
			PolicyNumber:      policyNumber,
			DriverID:          driverNumber,
			FirstName:         driverData["FirstName"].(string),
			LastName:          driverData["LastName"].(string),
			LicenseNumber:     driverData["LicenseNumber"].(string),
			LicenseState:      driverData["LicenseState"].(string),
			Age:               int(driverData["Age"].(float64)),
			DrivingExperience: int(driverData["DrivingExperience"].(float64)),
			Gender:            driverData["Gender"].(string),
			MaritalStatus:     driverData["MaritalStatus"].(string),
			DrivingCourse:     driverData["DrivingCourse"].(string),
			//MonthsSinceCourseCompleted: driverData["MonthsSinceCourseCompleted"].(nil),
			GoodStudent: driverData["GoodStudent"].(string),
			StudentAway: driverData["StudentAway"].(string),
		}

		drivers = append(drivers, d)
	}
	bundle.Drivers = drivers

	//  Vehicles data
	vehiclesData := data["VehicleForm"].([]interface{})
	var vehicles []models.Vehicle
	for _, vehicle := range vehiclesData {
		vehicleData := vehicle.(map[string]interface{})

		vehicleNumber := int(vehicleData["VehicleID"].(float64))

		coveragesData := vehicleData["Coverages"].([]interface{})

		var coverages []models.Coverage
		for _, coverage := range coveragesData {
			coverageData := coverage.(map[string]interface{})
			c := models.Coverage{
				EffectiveDate:  eff_date,
				CoverageCode:   coverageData["CoverageCode"].(string),
				CoverageOption: coverageData["CoverageOption"].(string),
			}
			coverages = append(coverages, c)
		}

		v := models.Vehicle{
			EffectiveDate:   eff_date,
			PolicyNumber:    policyNumber,
			VehicleID:       vehicleNumber,
			Year:            int(vehicleData["Year"].(float64)),
			Make:            vehicleData["Make"].(string),
			ModelCd:         vehicleData["Model"].(string),
			PrimaryUse:      vehicleData["PrimaryUse"].(string),
			VehicleType:     vehicleData["VehicleType"].(string),
			PrimaryOperator: vehicleData["PrimaryOperator"].(string),
			LoanORLease:     vehicleData["LoanORLease"].(string),
			Rideshare:       vehicleData["Rideshare"].(string),
			Coverages:       coverages,
		}
		vehicles = append(vehicles, v)
	}
	bundle.Vehicles = vehicles

	// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 			"message": "Internal server error",
	// 		})

	return nil

}

func ParseDate(DateValue string) time.Time {
	eff_date, err := time.Parse("2006-01-02T15:04:05.999Z", DateValue)
	if err != nil {
		logger.GetLogger().Info("Error parsing Date:", err)
		return time.Time{}
	}
	eff_date = time.Date(eff_date.Year(), eff_date.Month(), eff_date.Day(), 0, 0, 0, 0, eff_date.Location())
	return eff_date
}

func handleNullString(value interface{}) string {
	if value != nil {
		return value.(string)
	}
	return ""
}
