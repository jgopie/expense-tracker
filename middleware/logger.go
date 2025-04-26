package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logger(logger *log.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		stop := time.Now()
		latency := stop.Sub(start)
		logger.Printf("[%s] %s %d %s", c.Method(), c.Path(), c.Response().StatusCode(), latency)
		return err
	}
}
