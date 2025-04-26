package routes

import (
	"expense-tracker/controllers"
	"expense-tracker/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	api := app.Group("/api")
	users := api.Group("/users", middleware.Logger())
	users.Post("/register", controllers.Register)
	users.Post("/login", controllers.Login)
}
