package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/harishm11/Common/builder"
	"github.com/harishm11/Common/logger"
	models "github.com/harishm11/Common/models/PolicyModels"
)

const TRANSACTION_SERVICE_URL = "http://localhost:8003/transaction" 
func CallTransactionService(c *fiber.Ctx, transactionFormData *builder.TransactionFormData) (*models.Transaction, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	transactionRequestBody, err := json.Marshal(transactionFormData)
	if err != nil {
		logger.GetLogger().Error(err, "Error marshalling transaction data")
		return nil, errors.New("error marshalling transaction data")
	}

	if TRANSACTION_SERVICE_URL == "" {
		logger.GetLogger().Error(errors.New("TRANSACTION_SERVICE_URL not set"), "Transaction service URL not set")
		return nil, errors.New("transaction service URL not set")
	}

	resp, err := client.Post(TRANSACTION_SERVICE_URL, "application/json", bytes.NewBuffer(transactionRequestBody))
	if err != nil {
		logger.GetLogger().Error(err, "Failed to call transaction service")
		return nil, errors.New("failed to call transaction service")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		logger.GetLogger().Error(errors.New("failed to get valid response from transaction service"), "Transaction service response error")
		return nil, errors.New("transaction service response error")
	}

	var createdTransaction models.Transaction
	if err := json.NewDecoder(resp.Body).Decode(&createdTransaction); err != nil {
		logger.GetLogger().Error(err, "Error decoding transaction service response")
		return nil, errors.New("error decoding transaction service response")
	}

	return &createdTransaction, nil
}
