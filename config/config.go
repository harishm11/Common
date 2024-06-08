package config

import (
	"log"

	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("../../common/config/config.yaml")
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

}
func GetConfig(key string) string {
	return viper.GetString(key)
}
