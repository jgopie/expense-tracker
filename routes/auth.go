package routes

import (
	"expense-tracker/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	api := app.Group("/api")
	users := api.Group("/users")
	users.Post("/register", controllers.Register)
	users.Post("/login", controllers.Login)
	users.Get("/register", controllers.RenderRegistrationPage)
	users.Get("/login", controllers.RenderLoginPage)
}
