package routes

import (
	"expense-tracker/controllers"
	"expense-tracker/middleware"

	"github.com/gofiber/fiber/v2"
)

func ExpenseTrackerRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/transactions", controllers.GetTransactions, middleware.Protected())
	api.Post("/transactions", controllers.CreateTransaction, middleware.Protected())
	api.Delete("/transactions/:id", controllers.DeleteTransaction, middleware.Protected())
}
