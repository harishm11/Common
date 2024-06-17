package utils

import (
	"fmt"
	"reflect"
	"time"

	"github.com/gofiber/fiber/v2"
	models "github.com/harishm11/Common/models/PolicyModels"
)

const (
	NoChange models.ChangeIndicator = "N"
	Changed  models.ChangeIndicator = "C"
	Added    models.ChangeIndicator = "A"
	Deleted  models.ChangeIndicator = "D"
)

// FieldChangeIndicators represents the change indicators for a model
type FieldChangeIndicators map[string]models.ChangeIndicator

// CompareAndSetFieldIndicators compares two instances of a model and sets the change indicators

func CompareAndSetFieldIndicators(modelType reflect.Type, oldInstance, newInstance interface{}) (FieldChangeIndicators, error) {
	indicators := make(FieldChangeIndicators)

	if modelType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type is not a struct")
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldName := field.Name
		if fieldName != "FieldChangeIndicators" && fieldName != "RecordChangeIndicator" && fieldName != "Model" && fieldName != "Coverages" && fieldName != "ID" && fieldName != "EffectiveDate" {
			oldValue := reflect.ValueOf(oldInstance).FieldByName(fieldName)
			newValue := reflect.ValueOf(newInstance).FieldByName(fieldName)

			if oldValue.Kind() == reflect.Struct && oldValue.Type() == reflect.TypeOf(time.Time{}) {
				oldTime := oldValue.Interface().(time.Time).UTC() // Convert to UTC
				newTime := newValue.Interface().(time.Time).UTC() // Convert to UTC
				if oldTime != newTime {
					indicators[fieldName] = Changed
				} else {
					indicators[fieldName] = NoChange
				}
			} else {
				oldIsZero := reflect.DeepEqual(oldValue.Interface(), reflect.Zero(oldValue.Type()).Interface())
				newIsZero := reflect.DeepEqual(newValue.Interface(), reflect.Zero(newValue.Type()).Interface())

				if oldIsZero && !newIsZero {
					indicators[fieldName] = Added
				} else if !oldIsZero && newIsZero {
					indicators[fieldName] = Deleted
				} else if oldValue.Interface() != newValue.Interface() && oldIsZero && newIsZero {
					indicators[fieldName] = Changed
				} else {
					indicators[fieldName] = NoChange
				}
			}
		}
	}
	return indicators, nil
}

func SetIndicators(c *fiber.Ctx, bundle *models.Bundle, FromDbBundle *models.Bundle) error {

	// Compare and set indicators for Policy
	policyIndicators, err := CompareAndSetFieldIndicators(reflect.TypeOf(models.Policy{}), FromDbBundle.Policy, bundle.Policy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error"})
	}
	bundle.Policy.FieldChangeIndicators = models.ChangeIndicators(policyIndicators)
	SetRecordChangeIndicatorPolicy(bundle)

	// Compare and set indicators for PolicyHolder
	policyHolderIndicators, err := CompareAndSetFieldIndicators(reflect.TypeOf(models.PolicyHolder{}), FromDbBundle.PolicyHolder, bundle.PolicyHolder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error"})
	}
	bundle.PolicyHolder.FieldChangeIndicators = models.ChangeIndicators(policyHolderIndicators)
	SetRecordChangeIndicatorPolicyHolder(bundle)

	// Compare and set indicators for Current Carrier
	currentCarrierIndicators, err := CompareAndSetFieldIndicators(reflect.TypeOf(models.CurrentCarrierInfo{}), FromDbBundle.CurrentCarrier, bundle.CurrentCarrier)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error"})
	}
	bundle.CurrentCarrier.FieldChangeIndicators = models.ChangeIndicators(currentCarrierIndicators)
	SetRecordChangeIndicatorCurrentCarrier(bundle)

	// Iterate over each vehicle in the bundle
	for i, vehicle := range bundle.Vehicles {
		// Find the corresponding vehicle in FromDbBundle
		var dbVehicle *models.Vehicle
		for _, v := range FromDbBundle.Vehicles {
			if v.VehicleID == vehicle.VehicleID {
				dbVehicle = &v
				break
			}
		}
		if dbVehicle == nil {
			// Vehicle not found in FromDbBundle
			dbVehicle = &models.Vehicle{}
		}

		// Compare and set indicators for the current vehicle
		vehicleIndicators, err := CompareAndSetFieldIndicators(reflect.TypeOf(models.Vehicle{}), *dbVehicle, vehicle)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error"})
		}
		// Assign the change indicators to the current vehicle
		bundle.Vehicles[i].FieldChangeIndicators = models.ChangeIndicators(vehicleIndicators)

		// Now, update change indicators for each coverage in the vehicle
		for j, coverage := range vehicle.Coverages {
			if j < len(dbVehicle.Coverages) && coverage.CoverageCode == dbVehicle.Coverages[j].CoverageCode {
				coverageIndicators, err := CompareAndSetFieldIndicators(reflect.TypeOf(models.Coverage{}), dbVehicle.Coverages[j], coverage)
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error"})
				}
				// Assign the change indicators to the current coverage
				bundle.Vehicles[i].Coverages[j].FieldChangeIndicators = models.ChangeIndicators(coverageIndicators)
			}
		}
	}

	SetRecordChangeIndicatorVehicle(bundle, FromDbBundle)

	// Iterate over each driver in the bundle
	for i, driver := range bundle.Drivers {
		// Find the corresponding driver in FromDbBundle
		var dbDriver *models.Driver
		for _, d := range FromDbBundle.Drivers {
			if d.DriverID == driver.DriverID {
				dbDriver = &d
				break
			}
		}
		if dbDriver == nil {
			// Driver not found in FromDbBundle
			dbDriver = &models.Driver{}
		}

		// Compare and set indicators for the current driver
		driverIndicators, err := CompareAndSetFieldIndicators(reflect.TypeOf(models.Driver{}), *dbDriver, driver)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error"})
		}
		// Assign the change indicators to the current driver
		bundle.Drivers[i].FieldChangeIndicators = models.ChangeIndicators(driverIndicators)
	}

	SetRecordChangeIndicatorDriver(bundle, FromDbBundle)

	return nil

}
func SetRecordChangeIndicatorDriver(bundle1, bundle2 *models.Bundle) error {
	// Create a map to store drivers' IDs and their indices in bundle1
	driverIndices := make(map[int]int)

	// Populate the map for drivers in bundle1
	for i, driver1 := range bundle1.Drivers {
		driverIndices[driver1.DriverID] = i
	}

	// Iterate over drivers in bundle2 and compare with drivers in bundle1
	for _, driver2 := range bundle2.Drivers {
		// Check if driver2 exists in bundle1
		if idx, ok := driverIndices[driver2.DriverID]; ok {
			driver1 := bundle1.Drivers[idx]
			// Check if any field change indicators for this driver are 'C'
			if anyFieldChangeIndicatorsChanged(driver1.FieldChangeIndicators) {
				bundle1.Drivers[idx].RecordChangeIndicator = "C" // Set RecordChangeIndicator to 'C' if any field changes
			}
			// Remove the driver index from the map
			delete(driverIndices, driver2.DriverID)
		} else {
			// Driver in bundle2 but not in bundle1 (added)
			bundle1.Drivers = append(bundle1.Drivers, driver2)
			bundle1.Drivers[len(bundle1.Drivers)-1].RecordChangeIndicator = "D"
		}
	}

	// Set RecordChangeIndicator to 'D' for drivers in bundle1 but not in bundle2 (removed)
	for _, idx := range driverIndices {
		bundle1.Drivers[idx].RecordChangeIndicator = "A"
	}

	return nil
}

func anyFieldChangeIndicatorsChanged(indicators models.ChangeIndicators) bool {
	for _, indicator := range indicators {
		if indicator == "C" {
			return true
		}
	}
	return false
}

func allFieldChangeIndicatorsAdded(indicators models.ChangeIndicators) bool {
	for _, indicator := range indicators {
		if indicator != "A" {
			return false
		}
	}
	return true
}

func SetRecordChangeIndicatorVehicle(bundle1, bundle2 *models.Bundle) error {
	// Create a map to store vehicles' IDs and their indices in bundle1
	vehicleIndices := make(map[int]int)

	// Populate the map for vehicles in bundle1
	for i, vehicle1 := range bundle1.Vehicles {
		vehicleIndices[vehicle1.VehicleID] = i
	}

	// // If bundle2.Vehicles is empty, set RecordChangeIndicator to 'D' for all vehicles in bundle1
	// if len(bundle2.Vehicles) == 0 {
	//     for i := range bundle1.Vehicles {
	//         bundle1.Vehicles[i].RecordChangeIndicator = "A"
	//     }
	// }

	// Iterate over vehicles in bundle2 and compare with vehicles in bundle1
	for _, vehicle2 := range bundle2.Vehicles {
		// Check if vehicle2 exists in bundle1
		if idx, ok := vehicleIndices[vehicle2.VehicleID]; ok {
			vehicle1 := bundle1.Vehicles[idx]
			// Check if any field change indicators for this vehicle are 'C'
			if anyFieldChangeIndicatorsChanged(vehicle1.FieldChangeIndicators) {
				bundle1.Vehicles[idx].RecordChangeIndicator = "C" // Set RecordChangeIndicator to 'C' if any field changes
			}
			// Compare coverages for the current vehicle
			SetRecordChangeIndicatorCoverage(&bundle1.Vehicles[idx], &vehicle2)
			// Remove the vehicle index from the map
			delete(vehicleIndices, vehicle2.VehicleID)
		} else {
			// Vehicle in bundle2 but not in bundle1 (added)
			bundle1.Vehicles = append(bundle1.Vehicles, vehicle2)
			bundle1.Vehicles[len(bundle1.Vehicles)-1].RecordChangeIndicator = "D"
		}
	}

	// Set RecordChangeIndicator to 'D' for vehicles in bundle1 but not in bundle2 (removed)
	for _, idx := range vehicleIndices {
		bundle1.Vehicles[idx].RecordChangeIndicator = "A"
	}

	return nil
}

func SetRecordChangeIndicatorCoverage(vehicle1, vehicle2 *models.Vehicle) {
	// Create a map to store coverages' IDs and their indices in vehicle1
	coverageIndices := make(map[string]int)

	// Populate the map for coverages in vehicle1
	for i, coverage1 := range vehicle1.Coverages {
		coverageIndices[coverage1.CoverageCode] = i
	}

	// Iterate over coverages in vehicle2 and compare with coverages in vehicle1
	for _, coverage2 := range vehicle2.Coverages {
		// Check if coverage2 exists in vehicle1
		if idx, ok := coverageIndices[coverage2.CoverageCode]; ok {
			coverage1 := vehicle1.Coverages[idx]
			// Check if any field change indicators for this coverage are 'C'
			if anyFieldChangeIndicatorsChanged(coverage1.FieldChangeIndicators) {
				vehicle1.Coverages[idx].RecordChangeIndicator = "C" // Set RecordChangeIndicator to 'C' if any field changes
			}
			// Remove the coverage index from the map
			delete(coverageIndices, coverage2.CoverageCode)
		} else {
			// Coverage in vehicle2 but not in vehicle1 (added)
			vehicle1.Coverages = append(vehicle1.Coverages, coverage2)
			vehicle1.Coverages[len(vehicle1.Coverages)-1].RecordChangeIndicator = "D"
		}
	}

	// Set RecordChangeIndicator to 'D' for coverages in vehicle1 but not in vehicle2 (removed)
	for _, idx := range coverageIndices {
		vehicle1.Coverages[idx].RecordChangeIndicator = "A"
	}
}

func SetRecordChangeIndicatorPolicy(bundle *models.Bundle) error {
	// Get the policy from each bundle
	if anyFieldChangeIndicatorsChanged(bundle.Policy.FieldChangeIndicators) {
		bundle.Policy.RecordChangeIndicator = "C" // Set RecordChangeIndicator to 'C' if any field changes
	} else if allFieldChangeIndicatorsAdded(bundle.Policy.FieldChangeIndicators) {
		bundle.Policy.RecordChangeIndicator = "A"
	} else {
		// If the policies are not equal, set the record change indicator to "C" (changed)
		bundle.Policy.RecordChangeIndicator = "N"
	}

	return nil
}

func SetRecordChangeIndicatorPolicyHolder(bundle *models.Bundle) error {
	// Get the policy from each bundle
	if anyFieldChangeIndicatorsChanged(bundle.PolicyHolder.FieldChangeIndicators) {
		bundle.PolicyHolder.RecordChangeIndicator = "C" // Set RecordChangeIndicator to 'C' if any field changes
	} else if allFieldChangeIndicatorsAdded(bundle.PolicyHolder.FieldChangeIndicators) {
		bundle.PolicyHolder.RecordChangeIndicator = "A"
	} else {
		// If the policies are not equal, set the record change indicator to "C" (changed)
		bundle.PolicyHolder.RecordChangeIndicator = "N"
	}

	return nil
}

func SetRecordChangeIndicatorCurrentCarrier(bundle *models.Bundle) error {
	// Get the policy from each bundle
	if anyFieldChangeIndicatorsChanged(bundle.CurrentCarrier.FieldChangeIndicators) {
		bundle.CurrentCarrier.RecordChangeIndicator = "C" // Set RecordChangeIndicator to 'C' if any field changes
	} else if allFieldChangeIndicatorsAdded(bundle.CurrentCarrier.FieldChangeIndicators) {
		bundle.CurrentCarrier.RecordChangeIndicator = "A"
	} else {
		// If the policies are not equal, set the record change indicator to "C" (changed)
		bundle.CurrentCarrier.RecordChangeIndicator = "N"
	}

	return nil
}
