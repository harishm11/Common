package config

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
)

func SetupCors() fiber.Handler {
	env := viper.GetString("env")
	corsConfig := viper.Sub(env + ".cors")
	if corsConfig == nil {
		log.Fatalf("No CORS configuration found for environment: %s", env)
	}

	return cors.New(cors.Config{
		AllowHeaders:     corsConfig.GetString("allow_headers"),
		AllowOrigins:     corsConfig.GetString("allow_origins"),
		AllowCredentials: corsConfig.GetBool("allow_credentials"),
		AllowMethods:     corsConfig.GetString("allow_methods"),
	})
}
