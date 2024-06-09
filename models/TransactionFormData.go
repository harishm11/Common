// TransactionFormData represents the data required to create a new transaction
package models

import (
	"time"
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
