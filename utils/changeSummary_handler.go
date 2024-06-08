package utils

import (
	"encoding/json"
	"fmt"

	"github.com/harishm11/PolicyProcessor_V1.0/common/models"
	workflowmodels "github.com/harishm11/PolicyProcessor_V1.0/services/workflow_service/models"
)

// RecordChangeIndicators represents the change indicators for records
type RecordChangeIndicators map[string]models.ChangeIndicator

type FieldChange struct {
	Record   interface{}
	OldValue interface{}
	NewValue interface{}
}

// ChangeSummary represents the summary of changes
type ChangeSummary struct {
	Fields  FieldChangeIndicators  `json:"fields"`
	Records RecordChangeIndicators `json:"records"`
	Changes map[string]FieldChange
}

// StringToChangeIndicator converts a string to models.ChangeIndicator
func StringToChangeIndicator(str string) models.ChangeIndicator {
	switch str {
	case "N":
		return "No Change"
	case "C":
		return "Updated"
	case "A":
		return "Added"
	case "D":
		return "Deleted"
	default:
		return "No Change"
	}
}

// GenerateChangeSummary generates a summary of changes based on the provided indicators
// GenerateChangeSummary generates a summary of changes based on the provided indicators
func GenerateChangeSummary(currBundle, priorBundle *workflowmodels.Bundle) ([]byte, error) {
	summary := ChangeSummary{
		Fields:  make(FieldChangeIndicators),
		Records: make(RecordChangeIndicators),
		Changes: make(map[string]FieldChange),
	}

	// Helper function to capture changes
	captureChanges := func(field string, oldValue, newValue interface{}, indicator models.ChangeIndicator, recordkey string) {
		if indicator != "N" {
			summary.Changes[field] = FieldChange{
				Record:   recordkey,
				OldValue: oldValue,
				NewValue: newValue,
			}
		}
	}

	// Generate change summary for Policy
	for field, indicator := range currBundle.Policy.FieldChangeIndicators {
		summary.Fields[field] = indicator
		newValue := currBundle.Policy.GetFieldValue(field)
		oldValue := priorBundle.Policy.GetFieldValue(field)
		captureChanges(field, oldValue, newValue, indicator, "Policy")
	}
	summary.Records["Policy"] = StringToChangeIndicator(currBundle.Policy.RecordChangeIndicator)

	// // Generate change summary for PolicyHolder
	// for field, indicator := range currBundle.PolicyHolder.FieldChangeIndicators {
	// 	summary.Fields[field] = indicator
	// 	oldValue := currBundle.PolicyHolder.GetFieldValue(field)
	// 	newValue := priorBundle.PolicyHolder.GetFieldValue(field)
	// 	captureChanges(field, oldValue, newValue, indicator, "Policy Holder")
	// }
	// summary.Records["PolicyHolder"] = StringToChangeIndicator(currBundle.PolicyHolder.RecordChangeIndicator)

	// Generate change summary for CurrentCarrier
	for field, indicator := range currBundle.CurrentCarrier.FieldChangeIndicators {
		summary.Fields[field] = indicator
		newValue := currBundle.CurrentCarrier.GetFieldValue(field)
		oldValue := priorBundle.CurrentCarrier.GetFieldValue(field)
		captureChanges(field, oldValue, newValue, indicator, "Current Carrier")
	}
	summary.Records["CurrentCarrier"] = StringToChangeIndicator(currBundle.CurrentCarrier.RecordChangeIndicator)

	// Generate change summary for Vehicles
	for _, oldVehicle := range currBundle.Vehicles {
		recordKey := fmt.Sprintf(" %d  - %d %s %s   ", oldVehicle.VehicleID, oldVehicle.Year, oldVehicle.Make, oldVehicle.ModelCd)
		summary.Records[recordKey] = StringToChangeIndicator(oldVehicle.RecordChangeIndicator)
		if summary.Records[recordKey] != "Added" && summary.Records[recordKey] != "Deleted" && summary.Records[recordKey] != "No Change" {
			correspondingVeh := priorBundle.VehicleByID(oldVehicle.VehicleID)
			// Generate change summary for the current vehicle
			for field, indicator := range oldVehicle.FieldChangeIndicators {
				summary.Fields[field] = indicator
				oldValue := oldVehicle.GetFieldValue(field)
				newValue := oldVehicle.GetFieldValue(field)
				if correspondingVeh != nil {
					newValue = correspondingVeh.GetFieldValue(field)
				}
				captureChanges(field, oldValue, newValue, indicator, recordKey)
			}
			// Generate change summary for coverages
			for _, oldCoverage := range oldVehicle.Coverages {
				coverageKey := fmt.Sprintf("%s: %s", recordKey, oldCoverage.CoverageCode)
				summary.Records[coverageKey] = StringToChangeIndicator(oldCoverage.RecordChangeIndicator)

				// Generate change summary for the current coverage
				for field, indicator := range oldCoverage.FieldChangeIndicators {
					oldValue := oldCoverage.GetFieldValue(field)
					newValue := oldValue // Assume new value is the same as old value

					// Check if the corresponding coverage exists in the new bundle
					correspondingCoverage := correspondingVeh.CoverageByID(oldCoverage.CoverageCode)
					if correspondingCoverage != nil {
						// Get the new value if the corresponding coverage exists
						newValue = correspondingCoverage.GetFieldValue(field)
					}

					// Capture changes if the indicator is "C" and values are different
					captureChanges(field, oldValue, newValue, indicator, coverageKey)

				}
			}
		}

	}

	// Generate change summary for Drivers
	for _, oldDriver := range currBundle.Drivers {
		recordKey := fmt.Sprintf(" %d - %s %s   ", oldDriver.DriverID, oldDriver.FirstName, oldDriver.LastName)
		summary.Records[recordKey] = StringToChangeIndicator(oldDriver.RecordChangeIndicator)
		if summary.Records[recordKey] != "Added" && summary.Records[recordKey] != "Deleted" && summary.Records[recordKey] != "No Change" {
			// Generate change summary for the current driver
			for field, indicator := range oldDriver.FieldChangeIndicators {
				summary.Fields[field] = indicator
				oldValue := oldDriver.GetFieldValue(field)
				newValue := oldDriver.GetFieldValue(field)
				correspondingVeh := priorBundle.DriverByID(oldDriver.DriverID)
				if correspondingVeh != nil {
					newValue = correspondingVeh.GetFieldValue(field)
				}
				captureChanges(field, oldValue, newValue, indicator, recordKey)
			}
		}
	}

	// Marshal summary into JSON
	summaryJSON, err := json.Marshal(summary)
	if err != nil {
		return nil, err
	}

	return summaryJSON, nil
}
