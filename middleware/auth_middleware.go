package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("myjwtsecret")

// Middleware to check if user is authenticated
func AuthRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("token")

	if cookie == "" {
		return c.Redirect("/login")
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Redirect("/login")
	}

	return c.Next()
}

func CheckAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies("token")

		path := c.Path()
		if path == "/login" || path == "/register" {
			return c.Next()
		}

		if cookie == "" {
			return c.Redirect("/login")
		}

		token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Redirect("/login")
		}

		return c.Next()
	}
}
