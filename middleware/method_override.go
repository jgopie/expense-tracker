package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type MethodOverrideConfig struct {
	// The getter receives the method from the request
	Getter func(c *fiber.Ctx) string
}

var DefaultConfig = MethodOverrideConfig{Getter: methodFromForm}

func MethodOverride(config ...MethodOverrideConfig) fiber.Handler {
	// Set default config
	cfg := DefaultConfig

	// Override config if provided
	if len(config) > 0 {
		cfg = config[0]
		if cfg.Getter == nil {
			cfg.Getter = DefaultConfig.Getter
		}
	}

	return func(c *fiber.Ctx) error {
		// Only override POST methods
		if string(c.Request().Header.Method()) == fiber.MethodPost {
			method := cfg.Getter(c)
			if method != "" {
				method = strings.ToUpper(method)
				// Only allow certain HTTP methods
				switch method {
				case fiber.MethodPut, fiber.MethodPatch, fiber.MethodDelete:
					c.Request().Header.SetMethod(method)
				}
			}
		}

		return c.Next()
	}
}

func methodFromForm(c *fiber.Ctx) string {
	return c.FormValue("_method")
}

// Alternative method in case we want to check headers
// func methodFromHeader(c *fiber.Ctx) string {
// 	return c.Get("X-HTTP-Method-Override")
// }
