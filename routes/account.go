package routes

import (
	"expense-tracker/controllers"

	"github.com/gofiber/fiber/v2"
)

func AccountRoutes(app *fiber.App) {
	app.Get("/accounts", controllers.GetAccounts)
	app.Post("/accounts", controllers.CreateAccount)
	app.Get("/accounts/new", func(c *fiber.Ctx) error {
		return c.Render("new_account", fiber.Map{"Title": "New Account"})
	})
}
