package middleware

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logger() fiber.Handler {
	file, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to open log file %v", err)
	}
	logger := log.New(file, "", log.LstdFlags)
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		stop := time.Now()
		latency := stop.Sub(start)
		logger.Printf("[%s] %s %d %s", c.Method(), c.Path(), c.Response().StatusCode(), latency)
		return err
	}
}
