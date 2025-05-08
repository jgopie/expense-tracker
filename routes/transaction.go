package routes

import (
	"expense-tracker/controllers"

	"github.com/gofiber/fiber/v2"
)

func ExpenseTrackerRoutes(app *fiber.App) {
	app.Get("/transactions/add", controllers.RenderAddTransactionPage)
	app.Post("/transactions/add", controllers.CreateTransaction)
	app.Delete("/transactions/:id", controllers.DeleteTransaction)
	app.Post("/transactions/:id/delete", controllers.DeleteTransaction)
}
