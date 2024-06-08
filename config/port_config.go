// config.go

package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// GetPortForService retrieves the port number for a given service name from config.yaml
func GetPortForService(serviceName string) (int, error) {

	// Read the port number for the specified service name from config.yaml
	port := viper.GetInt(fmt.Sprintf("%s_port", serviceName))
	if port == 0 {
		return 0, fmt.Errorf("port not found for service: %s", serviceName)
	}

	return port, nil
}
