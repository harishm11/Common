package utils

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/harishm11/Common/config"
	"github.com/harishm11/Common/logger"
	models "github.com/harishm11/Common/models/PolicyModels"
)

func SaveRequest(c *fiber.Ctx, requestData map[string]interface{}, bundle *models.Bundle, mode string, workflowName string) error {
	policyProcessorDB := config.GetDBConn()

	eff_date := bundle.Policy.EffectiveDate
	accountNumber := bundle.Policy.AccountNumber
	policyNumber := bundle.Policy.PolicyNumber
	transactionNumber := bundle.Transaction.TransactionNumber
	lob := bundle.Policy.LineofBusiness
	lobForm := requestData["LobForm"].(map[string]interface{})
	lobForm["PolicyNumber"] = policyNumber
	lobForm["TransactionNumber"] = transactionNumber
	lobForm["TransactionType"] = bundle.Transaction.TransactionType

	jsonRequest, err := json.Marshal(requestData)
	if err != nil {
		logger.GetLogger().Error(err, "Error")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error"})
	}
	var status string
	switch mode {
	case "quote":
		status = "draft"
	case "save":
		status = "draft"
	case "bind":
		status = "bound"
	default:
		status = "failed"
	}

	logger.GetLogger().Info(status)
	request := models.Transaction{
		EffectiveDate:     eff_date,
		AccountNumber:     accountNumber,
		PolicyNumber:      policyNumber,
		TransactionNumber: transactionNumber,
		Status:            capitalizeFirstChar(status),
		TransactionType:   bundle.Transaction.TransactionType,
		Lob:               lob,
		JSONRequest:       string(jsonRequest),
		JSONResponse:      "",
	}

	// Check if a record with the same account and policy numbers exists
	var existingRequest models.Transaction
	if err := policyProcessorDB.Where("account_number = ? AND policy_number = ?  AND transaction_number = ?", accountNumber, policyNumber, transactionNumber).First(&existingRequest).Error; err != nil {
		// If no record found, create a new one
		if policyProcessorDB.Create(&request).Error != nil {
			logger.GetLogger().Error(err, "Error")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error creating RequestResponse"})
		}
	} else {
		// If record found, update the existing one
		existingRequest.JSONRequest = request.JSONRequest
		existingRequest.Status = request.Status
		if policyProcessorDB.Save(&existingRequest).Error != nil {
			logger.GetLogger().Error(err, "Error")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error updating RequestResponse"})
		}
	}

	return nil
}

func capitalizeFirstChar(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
