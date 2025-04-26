package routes

import (
	"expense-tracker/controllers"

	"github.com/gofiber/fiber/v2"
)

func ExpenseTrackerRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/transactions", controllers.GetTransactions)
	api.Post("/transactions", controllers.CreateTransaction)
	api.Delete("/transactions/:id", controllers.DeleteTransaction)
}
