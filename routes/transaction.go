package routes

import (
	"expense-tracker/controllers"

	"github.com/gofiber/fiber/v2"
)

func ExpenseTrackerRoutes(app *fiber.App) {
	app.Get("/transactions", controllers.GetTransactions)
	app.Post("/transactions", controllers.CreateTransaction)
	app.Delete("/transactions/:id", controllers.DeleteTransaction)
	app.Get("/add", controllers.RenderAddTransactionPage)
}
