package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/harishm11/API-Gateway/routes"
)

func main() {
	app := fiber.New()

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
