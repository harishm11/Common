package builder

import (
	"time"

	models "github.com/harishm11/Common/models/PolicyModels"
	ratingmodels "github.com/harishm11/Common/models/RatingModels"
)

func BuildRateRequest(bundle *models.Bundle) (*ratingmodels.RateRequest, error) {
	// Prepare the RateRequest from the bundle data
	// This is just an example. You need to map your bundle fields to the RateRequest fields.
	rateRequest := &ratingmodels.RateRequest{
		AccountNumber:     bundle.Policy.AccountNumber,
		PolicyNumber:      bundle.Policy.PolicyNumber,
		TransactionNumber: bundle.Transaction.TransactionNumber,
		//EffectiveDate:     bundle.Policy.EffectiveDate,
		//TransactionDate:   bundle.Transaction.TransactionDate,
		EffectiveDate:   time.Now(),
		TransactionDate: time.Now(),
		VehDetails:      make([]ratingmodels.ReqVehData, len(bundle.Vehicles)),
		DrvDetails:      make([]ratingmodels.ReqDrvData, len(bundle.Drivers)),
	}

	// Map VehDetails from bundle to RateRequest
	for i, veh := range bundle.Vehicles {
		rateRequest.VehDetails[i] = ratingmodels.ReqVehData{
			ID:              veh.ID,
			VIN:             veh.VIN,
			PolicyNumber:    veh.PolicyNumber,
			Year:            veh.Year,
			Make:            veh.Make,
			ModelCd:         veh.ModelCd,
			PrimaryUse:      veh.PrimaryUse,
			VehicleType:     veh.VehicleType,
			PrimaryOperator: veh.PrimaryOperator,
			LoanORLease:     veh.LoanORLease,
			Rideshare:       veh.Rideshare,
			CvgDetails:      make([]ratingmodels.ReqCvgData, len(veh.Coverages)),
		}

		// Map Coverages for each vehicle
		for j, cvg := range veh.Coverages {
			rateRequest.VehDetails[i].CvgDetails[j] = ratingmodels.ReqCvgData{
				ID:             cvg.ID,
				VehicleID:      veh.ID,
				CoverageCode:   cvg.CoverageCode,
				CvgSymbol:      cvg.CvgSymbol,
				CoverageOption: cvg.CoverageOption,
			}
		}
	}

	// Map DrvDetails from bundle to RateRequest
	for i, drv := range bundle.Drivers {
		rateRequest.DrvDetails[i] = ratingmodels.ReqDrvData{
			ID:                         drv.ID,
			PolicyNumber:               drv.PolicyNumber,
			FirstName:                  drv.FirstName,
			LastName:                   drv.LastName,
			LicenseNumber:              drv.LicenseNumber,
			LicenseState:               drv.LicenseState,
			Age:                        drv.Age,
			DrivingExperience:          drv.DrivingExperience,
			Gender:                     drv.Gender,
			MaritalStatus:              drv.MaritalStatus,
			DrivingCourse:              drv.DrivingCourse,
			MonthsSinceCourseCompleted: drv.MonthsSinceCourseCompleted,
			GoodStudent:                drv.GoodStudent,
			StudentAway:                drv.StudentAway,
		}
	}

	return rateRequest, nil
}
