package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	RatingServiceURL      string
	TransactionServiceURL string
	AccountServiceURL     string
	PolicyServiceURL      string
}

var AppConfig Config

func InitConfig() {
	viper.SetConfigName("..config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	viper.GetString("base_dir")
	env := viper.GetString("env")
	if env == "" {
		env = "development"
		viper.Set("env", env)
	}

	AppConfig.RatingServiceURL = viper.GetString("RATING_SERVICE_URL")
	AppConfig.TransactionServiceURL = viper.GetString("TRANSACTION_SERVICE_URL")
	AppConfig.PolicyServiceURL = viper.GetString("POLICY_SERVICE_URL")
	AppConfig.AccountServiceURL = viper.GetString("ACCOUNT_SERVICE_URL")

}
func GetConfig(key string) string {
	return viper.GetString(key)
}
