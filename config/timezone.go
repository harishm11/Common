package config

import (
	"log"
	"os"
	"time"
)

func SetupTimezone() {

	timezone := os.Getenv("TIME_ZONE")
	if timezone == "" {
		timezone = "UTC"
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Fatalf("Failed to load location: %s", err)
	}

	time.Local = loc
}
