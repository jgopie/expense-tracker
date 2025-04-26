package routes

import (
	"expense-tracker/controllers"
	"expense-tracker/middleware"

	"github.com/gofiber/fiber/v2"
)

func ExpenseTrackerRoutes(app *fiber.App) {
	api := app.Group("/api")
	transactions := api.Group("/transactions", middleware.Protected(), middleware.Logger())
	transactions.Get("/", controllers.GetTransactions)
	transactions.Post("/", controllers.CreateTransaction)
	transactions.Delete("/:id", controllers.DeleteTransaction)
}
