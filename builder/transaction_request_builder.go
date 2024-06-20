package builder

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/harishm11/Common/logger"
	models "github.com/harishm11/Common/models/PolicyModels"
)

type TransactionFormData struct {
	AccountNumber     int       `json:"accountNumber"`
	PolicyNumber      int       `json:"policyNumber"`
	TransactionNumber int       `json:"transactionNumber"`
	TransactionType   string    `json:"transactionType"`
	Lob               string    `json:"lob"`
	EffectiveDate     time.Time `json:"effectiveDate"`
	TransactionDate   time.Time `json:"transactionDate"`
	Status            string    `json:"status"`
	JSONRequest       string    `json:"jsonRequest"`
	JSONResponse      string    `json:"jsonResponse"`
}

func BuildTransactionRequest(requestData map[string]interface{}, bundle *models.Bundle, mode string) (TransactionFormData, error) {
	effDate := bundle.Policy.EffectiveDate
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
		return TransactionFormData{}, errors.New("error marshalling request data")
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

	transactionFormData := TransactionFormData{
		AccountNumber:     accountNumber,
		PolicyNumber:      policyNumber,
		TransactionNumber: transactionNumber,
		TransactionType:   bundle.Transaction.TransactionType,
		Lob:               lob,
		EffectiveDate:     effDate,
		TransactionDate:   bundle.Transaction.TransactionDate,
		Status:            capitalizeFirstChar(status),
		JSONRequest:       string(jsonRequest),
		JSONResponse:      "",
	}

	return transactionFormData, nil
}

func capitalizeFirstChar(str string) string {
	if len(str) == 0 {
		return str
	}
	return string(str[0]-32) + str[1:]
}
