package config

import (
	"fmt"
	"os"
	"strconv"
)

// GetPortForService retrieves the port number for the specified service from the environment variables.
func GetPortForService(serviceName string) (int, error) {
	var portStr string

	// Use switch case to determine which environment variable to use
	switch serviceName {
	case "AccountService":
		portStr = os.Getenv("SERVICE_PORT_ACCOUNT")
	case "PolicyService":
		portStr = os.Getenv("SERVICE_PORT_POLICY")
	case "TransactionService":
		portStr = os.Getenv("SERVICE_PORT_TRANSACTION")
	case "TransactionService":
		portStr = os.Getenv("SERVICE_PORT_WORKFLOW")
	case "RatingService":
		portStr = os.Getenv("SERVICE_PORT_RATING")
	default:
		return 0, fmt.Errorf("no port configuration found for service: %s", serviceName)
	}

	if portStr == "" {
		return 0, fmt.Errorf("no port found for service: %s", serviceName)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("invalid port format for service: %s", serviceName)
	}

	return port, nil
}
