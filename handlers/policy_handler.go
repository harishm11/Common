package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/harishm11/Common/builder"
	"github.com/harishm11/Common/clients"
	"github.com/harishm11/Common/logger"
	models "github.com/harishm11/Common/models/PolicyModels"
)

func PolicyRequestHandler(bundle *models.Bundle, c *fiber.Ctx) error {

	policyData := builder.PolicyData{
		Policy:         bundle.Policy,
		CurrentCarrier: bundle.CurrentCarrier,
		PolicyHolder:   bundle.PolicyHolder,
		PolicyAddress:  bundle.PolicyAddress,
		Drivers:        bundle.Drivers,
		Vehicles:       bundle.Vehicles,
	}

	response, err := clients.CallPolicyService(policyData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	logger.GetLogger().Info("Policy created successfully")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Policy created successfully", "response": response})
}
