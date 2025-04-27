package main

import (
	"expense-tracker/config"
	"expense-tracker/middleware"
	"expense-tracker/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.ConnectDB()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))
	file, err := os.OpenFile("/root/logs/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to open log file %v", err)
	}
	defer file.Close()
	logger := log.New(file, "", log.LstdFlags)
	app.Use(middleware.Logger(logger))
	routes.AuthRoutes(app)
	routes.ExpenseTrackerRoutes(app)
	app.Listen("0.0.0.0:3000")
}
