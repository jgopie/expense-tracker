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
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if userID, ok := claims["user_id"].(float64); ok {
			c.Locals("user_id", userID)
		}
	}

	return c.Next()
}

func CheckAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip auth for these paths
		path := c.Path()
		if path == "/login" || path == "/register" {
			return c.Next()
		}

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

		// Extract claims and set user_id in locals
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userID, ok := claims["user_id"].(float64); ok {
				c.Locals("user_id", userID)
			} else {
				return c.Redirect("/login")
			}
		} else {
			return c.Redirect("/login")
		}

		return c.Next()
	}
}
