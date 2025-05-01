package main

import (
	"expense-tracker/config"
	"expense-tracker/middleware"
	"expense-tracker/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	config.ConnectDB()
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})
	app.Static("/static", "./static")

	// Flag explanation (just in case)
	// O_Append - New information is appended instead of overwriting the contents
	// O_Create - File is created if it doesn't exist
	// O_Wronlg - Write Only. File is opened for writing only, not reading
	// 0666 - Permissions, very permissive which is the norm for log files
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		logPath = "./logs/server.logs"
	}
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to open log file %v", err)
	}
	defer file.Close()
	logger := log.New(file, "", log.LstdFlags)
	app.Use(middleware.Logger(logger))
	app.Use(middleware.CheckAuth())
	routes.AuthRoutes(app)
	routes.ExpenseTrackerRoutes(app)
	app.Listen("0.0.0.0:3000")
}
