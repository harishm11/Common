package config

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SetupCors configures CORS middleware for the application.
func SetupCors() fiber.Handler {
	allowHeaders := os.Getenv("CORS_ALLOW_HEADERS")
	if allowHeaders == "" {
		allowHeaders = "Origin,Content-Type,Accept,Authorization"
	}

	allowOrigins := os.Getenv("CORS_ALLOW_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = "*"
	}

	return cors.New(cors.Config{
		AllowHeaders: allowHeaders,
		AllowOrigins: allowOrigins,
	})
}
