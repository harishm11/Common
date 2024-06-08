package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

func SetupTimezone() {
	env := viper.GetString("env")
	timezone := viper.GetString(env + ".timezone")
	if timezone == "" {
		timezone = "UTC"
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Fatalf("Failed to load location: %s", err)
	}

	time.Local = loc
}
