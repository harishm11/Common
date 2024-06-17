package utils

import (
	"github.com/gofiber/fiber/v2"

	"github.com/harishm11/Common/config"
	"github.com/harishm11/Common/logger"
	models "github.com/harishm11/Common/models/PolicyModels"
	"gorm.io/gorm"
)

func Push(c *fiber.Ctx, bundle *models.Bundle) (string, error) {
	policyProcessorDB, err := config.InitDatabase("PolicyProcessorDB")
	if err != nil {
		logger.GetLogger().Error(err, "Failed to initialize PolicyProcessorDB")
	}

	eff_date := bundle.Transaction.EffectiveDate
	updates := map[string]interface{}{
		"deleted_at": eff_date,
	}

	// Save the Coverage models first and set their foreign key (e.g., VehicleID)
	for i := range bundle.Vehicles {
		for j := range bundle.Vehicles[i].Coverages {
			if bundle.Vehicles[i].Coverages[j].RecordChangeIndicator == "C" || bundle.Vehicles[i].Coverages[j].RecordChangeIndicator == "A" {
				// Update or create the coverage
				if err := updateOrCreate(policyProcessorDB, &bundle.Vehicles[i].Coverages[j], map[string]interface{}{"id": bundle.Vehicles[i].Coverages[j].ID}, &bundle.Vehicles[i].Coverages[j]); err != nil {
					return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "update/create coverage"})
				}
			} else if bundle.Vehicles[i].Coverages[j].RecordChangeIndicator == "D" {
				// Delete the coverage
				if err := policyProcessorDB.Model(&bundle.Vehicles[i].Coverages[j]).Updates(updates).Error; err != nil {
					return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "delete coverage"})
				}
			}
		}
	}

	// Save the Vehicle models and set their foreign key (e.g., BundleID)
	for i := range bundle.Vehicles {
		// Set the foreign key for Coverages inside each Vehicle
		for j := range bundle.Vehicles[i].Coverages {
			bundle.Vehicles[i].Coverages[j].VehicleID = bundle.Vehicles[i].ID
		}
		if bundle.Vehicles[i].RecordChangeIndicator == "C" || bundle.Vehicles[i].RecordChangeIndicator == "A" {
			if err := updateOrCreate(policyProcessorDB, &bundle.Vehicles[i], map[string]interface{}{"id": bundle.Vehicles[i].ID}, &bundle.Vehicles[i]); err != nil {
				return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "update/create vehicle"})
			}
		} else if bundle.Vehicles[i].RecordChangeIndicator == "D" {
			// Delete the vehicle
			if err := policyProcessorDB.Model(&bundle.Vehicles[i]).Updates(updates).Error; err != nil {
				return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "delete vehicle"})
			}
		}
	}

	// Save the Policy model and set its foreign key (e.g., BundleID)
	if bundle.Policy.RecordChangeIndicator == "C" || bundle.Policy.RecordChangeIndicator == "A" {
		if err := updateOrCreate(policyProcessorDB, &bundle.Policy, map[string]interface{}{"id": bundle.Policy.ID}, &bundle.Policy); err != nil {
			return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "update/create policy"})
		}
	} else if bundle.Policy.RecordChangeIndicator == "D" {
		// Delete the policy

		if err := policyProcessorDB.Model(&bundle.Policy).Updates(updates).Error; err != nil {
			return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "delete policy"})
		}
	}

	// Save the Current carrier model and set its foreign key (e.g., BundleID)

	if bundle.CurrentCarrier.RecordChangeIndicator == "C" || bundle.CurrentCarrier.RecordChangeIndicator == "A" {
		if err := updateOrCreate(policyProcessorDB, &bundle.CurrentCarrier, map[string]interface{}{"id": bundle.CurrentCarrier.ID}, &bundle.CurrentCarrier); err != nil {
			return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "update/create current carrier"})
		}
	} else if bundle.CurrentCarrier.RecordChangeIndicator == "D" {
		// Delete the CurrentCarrier

		if err := policyProcessorDB.Model(&bundle.CurrentCarrier).Updates(updates).Error; err != nil {
			return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "delete current carrier"})
		}
	}

	// Save the Policy holder model and set its foreign key (e.g., BundleID)

	if bundle.PolicyHolder.RecordChangeIndicator == "C" || bundle.PolicyHolder.RecordChangeIndicator == "A" {
		if err := updateOrCreate(policyProcessorDB, &bundle.PolicyHolder, map[string]interface{}{"id": bundle.PolicyHolder.ID}, &bundle.PolicyHolder); err != nil {
			return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "update/create policyholder"})
		}
	} else if bundle.PolicyHolder.RecordChangeIndicator == "D" {
		// Delete the PolicyHolder

		if err := policyProcessorDB.Model(&bundle.PolicyHolder).Updates(updates).Error; err != nil {
			return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "delete policyholder"})
		}
	}

	if bundle.PolicyAddress.RecordChangeIndicator == "C" || bundle.PolicyAddress.RecordChangeIndicator == "A" {
		if err := updateOrCreate(policyProcessorDB, &bundle.PolicyAddress, map[string]interface{}{"id": bundle.PolicyAddress.ID}, &bundle.PolicyAddress); err != nil {
			return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "update/create policyaddress"})
		}
	} else if bundle.PolicyAddress.RecordChangeIndicator == "D" {
		// Delete the PolicyAddress

		if err := policyProcessorDB.Model(&bundle.PolicyAddress).Updates(updates).Error; err != nil {
			return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "delete policyaddress"})
		}
	}

	// Save the Transaction model and set its foreign key (e.g., BundleID)
	if err := updateOrCreate(policyProcessorDB, &bundle.Transaction, map[string]interface{}{"id": bundle.Transaction.ID}, &bundle.Transaction); err != nil {
		return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "update/create transaction"})
	}

	// Save the Drivers models and set their foreign key (e.g., BundleID)
	for i := range bundle.Drivers {

		if bundle.Drivers[i].RecordChangeIndicator == "C" || bundle.Drivers[i].RecordChangeIndicator == "A" {
			if err := updateOrCreate(policyProcessorDB, &bundle.Drivers[i], map[string]interface{}{"id": bundle.Drivers[i].ID}, &bundle.Drivers[i]); err != nil {
				return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "update/create driver"})
			}
		} else if bundle.Drivers[i].RecordChangeIndicator == "D" {
			// Delete the vehicle

			if err := policyProcessorDB.Model(&bundle.Drivers[i]).Updates(updates).Error; err != nil {
				return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "delete driver"})
			}
		}
	}

	return "Success", nil
}

func SubmissionPush(c *fiber.Ctx, bundle *models.Bundle) (string, error) {
	policyProcessorDB, err := config.InitDatabase("PolicyProcessorDB")
	if err != nil {
		logger.GetLogger().Error(err, "Failed to initialize PolicyProcessorDB")
	}

	// Save the Coverage models first and set their foreign key (e.g., VehicleID)
	for i := range bundle.Vehicles {
		for j := range bundle.Vehicles[i].Coverages {
			//  create the coverage
			if err := updateOrCreate(policyProcessorDB, &bundle.Vehicles[i].Coverages[j], map[string]interface{}{"id": bundle.Vehicles[i].Coverages[j].ID}, &bundle.Vehicles[i].Coverages[j]); err != nil {
				return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "update/create coverage"})
			}

		}
	}

	// Save the Vehicle models and set their foreign key (e.g., BundleID)
	for i := range bundle.Vehicles {
		// Set the foreign key for Coverages inside each Vehicle
		for j := range bundle.Vehicles[i].Coverages {
			bundle.Vehicles[i].Coverages[j].VehicleID = bundle.Vehicles[i].ID
		}
		if err := updateOrCreate(policyProcessorDB, &bundle.Vehicles[i], map[string]interface{}{"id": bundle.Vehicles[i].ID}, &bundle.Vehicles[i]); err != nil {
			return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "update/create vehicle"})
		}

	}

	// Save the Policy model and set its foreign key (e.g., BundleID)
	if err := updateOrCreate(policyProcessorDB, &bundle.Policy, map[string]interface{}{"id": bundle.Policy.ID}, &bundle.Policy); err != nil {
		return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "update/create policy"})
	}

	// Save the Current carrier model and set its foreign key (e.g., BundleID)

	if err := updateOrCreate(policyProcessorDB, &bundle.CurrentCarrier, map[string]interface{}{"id": bundle.CurrentCarrier.ID}, &bundle.CurrentCarrier); err != nil {
		return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "update/create current carrier"})
	}

	// Save the Policy holder model and set its foreign key (e.g., BundleID)

	if err := updateOrCreate(policyProcessorDB, &bundle.PolicyHolder, map[string]interface{}{"id": bundle.PolicyHolder.ID}, &bundle.PolicyHolder); err != nil {
		return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "update/create policyholder"})
	}
	// Save the Policy holder model and set its foreign key (e.g., BundleID)

	if err := updateOrCreate(policyProcessorDB, &bundle.PolicyAddress, map[string]interface{}{"id": bundle.PolicyAddress.ID}, &bundle.PolicyAddress); err != nil {
		return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "update/create policyaddress"})
	}
	// Save the Transaction model and set its foreign key (e.g., BundleID)
	if err := updateOrCreate(policyProcessorDB, &bundle.Transaction, map[string]interface{}{"id": bundle.Transaction.ID}, &bundle.Transaction); err != nil {
		return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "update/create transaction"})
	}

	// Save the Drivers models and set their foreign key (e.g., BundleID)
	for i := range bundle.Drivers {

		if err := updateOrCreate(policyProcessorDB, &bundle.Drivers[i], map[string]interface{}{"id": bundle.Drivers[i].ID}, &bundle.Drivers[i]); err != nil {
			return "Error", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "update/create driver"})
		}

	}

	return "Success", nil
}

func updateOrCreate(policyProcessorDB *gorm.DB, model interface{}, conditions interface{}, update interface{}) error {
	result := policyProcessorDB.Where(conditions).Assign(update).FirstOrCreate(model)

	if err := result.Error; err != nil {
		return err
	}
	return nil
}
