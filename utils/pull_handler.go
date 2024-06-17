package utils

import (
	"time"

	"github.com/harishm11/Common/config"
	"github.com/harishm11/Common/logger"
	models "github.com/harishm11/Common/models/PolicyModels"
)

func ShortPull(policyNum int, effectiveDate time.Time, tables []string) (*models.Bundle, error) {
	policyProcessorDB := config.GetDBConn()

	// Load the Bundle model from the database
	var bundle models.Bundle

	// // Load the Transaction model for the Bundle with policyNum and effectiveDate filters
	// transaction, err := transactionmodels.PullTransaction(policyProcessorDB, policyNum, effectiveDate)
	// if err != nil {
	// 	return nil, PullError(err, "Transaction")
	// }
	// bundle.Transaction = *transaction

	// Load the Policy model for the Bundle with policyNum and effectiveDate filters
	policy, err := models.PullPolicy(policyProcessorDB, policyNum, effectiveDate)
	if err != nil {
		return nil, PullError(err, "Policy")
	}
	bundle.Policy = *policy

	// Load the PolicyHolder model for the Bundle with policyNum and effectiveDate filters
	policyHolder, err := models.PullPolicyHolder(policyProcessorDB, policyNum, effectiveDate)
	if err != nil {
		return nil, PullError(err, "PolicyHolder")
	}
	bundle.PolicyHolder = *policyHolder

	// Load the PolicyAddress model for the Bundle with policyNum and effectiveDate filters
	policyAddress, err := models.PullPolicyAddress(policyProcessorDB, policyNum, effectiveDate)
	if err != nil {
		return nil, PullError(err, "PolicyAddress")
	}
	bundle.PolicyAddress = *policyAddress

	// Load the Current Carrier model for the Bundle with policyNum and effectiveDate filters
	currentCarrier, err := models.PullCurrentCarrier(policyProcessorDB, policyNum, effectiveDate)
	if err != nil {
		return nil, PullError(err, "CurrentCarrier")
	}
	bundle.CurrentCarrier = *currentCarrier

	// Load the Drivers models for the Bundle with policyNum and effectiveDate filters
	drivers, err := models.PullDrivers(policyProcessorDB, policyNum, effectiveDate)
	if err != nil {
		return nil, PullError(err, "Drivers")
	}
	bundle.Drivers = *drivers

	// Load the Vehicles models for the Bundle with policyNum and effectiveDate filters
	vehicles, err := models.PullVehicles(policyProcessorDB, policyNum, effectiveDate)
	if err != nil {
		return nil, PullError(err, "Vehicles")
	}
	bundle.Vehicles = *vehicles

	// Load the Coverages models for each Vehicle with policyNum and effectiveDate filters
	for i := range bundle.Vehicles {
		coverages, err := models.PullCoveragesForVehicle(policyProcessorDB, bundle.Vehicles[i].ID, effectiveDate)
		if err != nil {
			return nil, PullError(err, "Coverages")
		}
		bundle.Vehicles[i].Coverages = *coverages
	}

	return &bundle, nil
}

func PullError(err error, modelname string) error {
	if err != nil {
		logger.GetLogger().Info("Error pulling model:", err, modelname)
	}
	return err
}
