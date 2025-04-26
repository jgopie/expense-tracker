package main

import (
	"expense-tracker/config"
	"expense-tracker/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.ConnectDB()
	app := fiber.New()
	routes.ExpenseTrackerRoutes(app)
	app.Listen(":3000")
}
