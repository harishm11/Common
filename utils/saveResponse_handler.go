package utils

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"

	"github.com/harishm11/PolicyProcessor_V1.0/common/config"
	"github.com/harishm11/PolicyProcessor_V1.0/common/logger"
	ratingmodels "github.com/harishm11/PolicyProcessor_V1.0/services/rating_service/models"
	transactionmodels "github.com/harishm11/PolicyProcessor_V1.0/services/transaction_service/models"
)

func SaveResponse(c *fiber.Ctx, taskResponse interface{}, mode string) error {

	policyProcessorDB, err := config.InitDatabase("PolicyProcessorDB")
	if err != nil {
		logger.GetLogger().Error(err, "Failed to initialize PolicyProcessorDB")
	}
	var status string
	switch mode {
	case "quote":
		status = "quoted"
	case "save":
		status = "saved"
	case "bind":
		status = "bound"
	default:
		status = "failed"
	}
	jsonResponse, err := json.Marshal(taskResponse)
	if err != nil {
		logger.GetLogger().Error(err, "Error")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error"})
	}

	var rateResponse []ratingmodels.RateResponse
	var ok bool
	rateResponse, ok = taskResponse.([]ratingmodels.RateResponse)
	if !ok {
		logger.GetLogger().LogAnything(ok)
	}

	transactionNumber := rateResponse[0].TransactionNumber
	accountNumber := rateResponse[0].AccountNumber
	policyNumber := rateResponse[0].PolicyNumber
	var existingRequestResponse transactionmodels.Transaction
	if err := policyProcessorDB.Where("transaction_number = ? AND account_number = ? AND policy_number = ?", transactionNumber, accountNumber, policyNumber).First(&existingRequestResponse).Error; err != nil {
		logger.GetLogger().Error(err, "Error")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error"})
	}

	existingRequestResponse.JSONResponse = string(jsonResponse)
	existingRequestResponse.Status = capitalizeFirstChar(status)

	if err := policyProcessorDB.Save(&existingRequestResponse).Error; err != nil {
		logger.GetLogger().Error(err, "Error")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error"})
	}

	return nil
}
