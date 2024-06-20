package builder

import (
	"encoding/json"

	"github.com/harishm11/Common/logger"
	models "github.com/harishm11/Common/models/PolicyModels"
)

type PolicyData struct {
	Policy         models.Policy             `json:"policy"`
	CurrentCarrier models.CurrentCarrierInfo `json:"currentCarrier"`
	PolicyHolder   models.PolicyHolder       `json:"policyHolder"`
	PolicyAddress  models.PolicyAddress      `json:"policyAddress"`
	Drivers        []models.Driver           `json:"drivers"`
	Vehicles       []models.Vehicle          `json:"vehicles"`
}

func PreparePolicyRequest(policyData PolicyData) ([]byte, error) {
	requestBody, err := json.Marshal(policyData)
	if err != nil {
		logger.GetLogger().Error(err, "Error marshalling policy data")
		return nil, err
	}
	return requestBody, nil
}
