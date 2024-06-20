package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/harishm11/Common/builder"
	"github.com/harishm11/Common/clients"
	"github.com/harishm11/Common/logger"
	models "github.com/harishm11/Common/models/PolicyModels"
)

func TransactionRequestHandler(c *fiber.Ctx, requestData map[string]interface{}, bundle *models.Bundle, mode string, workflowName string) error {
	transactionFormData, err := builder.BuildTransactionRequest(requestData, bundle, mode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	createdTransaction, err := clients.CallTransactionService(c, &transactionFormData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	bundle.Transaction.TransactionNumber = createdTransaction.TransactionNumber

	logger.GetLogger().Info("Transaction created successfully: ", createdTransaction)
	return c.JSON(fiber.Map{"message": "Transaction created successfully", "transaction": createdTransaction})
}
