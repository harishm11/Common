package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/harishm11/Common/builder"
	"github.com/harishm11/Common/clients"
	"github.com/harishm11/Common/logger"
)

func PolicyRequestHandler(c *fiber.Ctx) error {
	var policyData builder.PolicyData

	if err := c.BodyParser(&policyData); err != nil {
		logger.GetLogger().Error(err, "Error parsing policy data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid policy data"})
	}

	response, err := clients.CallPolicyService(policyData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	logger.GetLogger().Info("Policy created successfully")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Policy created successfully", "response": response})
}
