package httpclients

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/harishm11/Common/logger"
	models "github.com/harishm11/Common/models/PolicyModels"
)

func SaveTransaction(c *fiber.Ctx, requestData map[string]interface{}, bundle *models.Bundle, mode string, workflowName string) error {
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
		logger.GetLogger().Error(err, "Error marshalling request data")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error marshalling request data"})
	}

	var status string
	switch mode {
	case "quote", "save":
		status = "draft"
	case "bind":
		status = "bound"
	default:
		status = "failed"
	}

	transactionFormData := models.TransactionFormData{
		AccountNumber:     accountNumber,
		PolicyNumber:      policyNumber,
		TransactionNumber: transactionNumber,
		TransactionType:   bundle.Transaction.TransactionType,
		Lob:               lob,
		EffectiveDate:     eff_date,
		TransactionDate:   bundle.Transaction.TransactionDate,
		Status:            capitalizeFirstChar(status),
		JSONRequest:       string(jsonRequest),
		JSONResponse:      "",
	}

	// Marshal transactionFormData to JSON
	transactionRequestBody, err := json.Marshal(transactionFormData)
	if err != nil {
		logger.GetLogger().Error(err, "Error marshalling transaction data")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error marshalling transaction data"})
	}

	TRANSACTION_SERVICE_URL := "http://localhost:8003/transaction"
	transactionServiceURL := TRANSACTION_SERVICE_URL

	//transactionServiceURL := os.Getenv("TRANSACTION_SERVICE_URL")
	if transactionServiceURL == "" {
		logger.GetLogger().Error(errors.New("TRANSACTION_SERVICE_URL not set"), "Transaction service URL not set")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Transaction service URL not set"})
	}

	// Make the HTTP POST request to the transaction service
	resp, err := http.Post(transactionServiceURL, "application/json", bytes.NewBuffer(transactionRequestBody))
	if err != nil {
		logger.GetLogger().Error(err, "Failed to call transaction service")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to call transaction service"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		logger.GetLogger().Error(errors.New("failed to get valid response from transaction service"), "Transaction service response error")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Transaction service response error"})
	}

	// Decode the response body
	var createdTransaction models.Transaction
	if err := json.NewDecoder(resp.Body).Decode(&createdTransaction); err != nil {
		logger.GetLogger().Error(err, "Error decoding transaction service response")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error decoding transaction service response"})
	}

	// update the transction number in bundle
	bundle.Transaction.TransactionNumber = createdTransaction.TransactionNumber

	logger.GetLogger().Info("Transaction created successfully: ", createdTransaction)
	return c.JSON(fiber.Map{"message": "Transaction created successfully", "transaction": createdTransaction})
}

func capitalizeFirstChar(str string) string {
	if len(str) == 0 {
		return str
	}
	return string(str[0]-32) + str[1:]
}
