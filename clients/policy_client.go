package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/harishm11/Common/builder"
	"github.com/harishm11/Common/logger"
)

func CallPolicyService(policyData builder.PolicyData, POLICY_SERVICE_URL string) (interface{}, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	requestBody, err := builder.PreparePolicyRequest(policyData)
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(POLICY_SERVICE_URL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		logger.GetLogger().Error(err, "Failed to call policy service")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		logger.GetLogger().Error(errors.New("failed to get valid response from policy service"), "Policy service response error")
		return nil, errors.New("policy service response error")
	}

	var response interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logger.GetLogger().Error(err, "Error decoding policy service response")
		return nil, err
	}

	return response, nil
}
