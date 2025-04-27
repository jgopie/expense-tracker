package routes

import (
	"expense-tracker/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	app.Get("/register", controllers.RenderRegistrationPage)
	app.Get("/login", controllers.RenderLoginPage)
	app.Post("/register", controllers.ProcessRegisterForm)
	app.Post("/login", controllers.ProcessLoginForm)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Expense Tracker",
		})
	})
}
