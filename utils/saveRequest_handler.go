package utils

import (
	"encoding/json"
	"strings"

	"github.com/harishm11/PolicyProcessor_V1.0/common/config"
	"github.com/harishm11/PolicyProcessor_V1.0/common/logger"
	workflowmodels "github.com/harishm11/PolicyProcessor_V1.0/services/workflow_service/models"

	"github.com/gofiber/fiber/v2"
	transactionmodels "github.com/harishm11/PolicyProcessor_V1.0/services/transaction_service/models"
)

func SaveRequest(c *fiber.Ctx, requestData map[string]interface{}, bundle *workflowmodels.Bundle, mode string, workflowName string) error {
	policyProcessorDB, err := config.InitDatabase("PolicyProcessorDB")
	if err != nil {
		logger.GetLogger().Error(err, "Failed to initialize PolicyProcessorDB")
	}

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
	request := transactionmodels.Transaction{
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
	var existingRequest transactionmodels.Transaction
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
