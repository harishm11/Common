package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/harishm11/API-Gateway/config"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/")

	api.Get("/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "API Gateway is running",
		})
	})

	api.All("/Rating/*", proxy.Balancer(proxy.Config{
		Servers: []string{config.AppConfig.RatingServiceURL},
	}))

	api.All("/Transaction/*", proxy.Balancer(proxy.Config{
		Servers: []string{config.AppConfig.TransactionServiceURL},
	}))

	api.All("/Policy/*", proxy.Balancer(proxy.Config{
		Servers: []string{config.AppConfig.PolicyServiceURL},
	}))

	api.All("/Account/*", proxy.Balancer(proxy.Config{
		Servers: []string{config.AppConfig.PolicyServiceURL},
	}))
}
