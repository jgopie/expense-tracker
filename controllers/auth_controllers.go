package controllers

import (
	"expense-tracker/config"
	"expense-tracker/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("myjwtsecret")

func RenderLoginPage(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Title": "Login - Expense Tracker",
	})
}

func RenderRegistrationPage(c *fiber.Ctx) error {
	return c.Render("registration", fiber.Map{
		"Title": "Registration - Expense Tracker",
	})
}

// HTML form handlers
func ProcessLoginForm(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return c.Render("login", fiber.Map{
			"Title": "Login - Expense Tracker",
			"Error": "Invalid email or password",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.Render("login", fiber.Map{
			"Title": "Login - Expense Tracker",
			"Error": "Invalid email or password",
		})
	}

	// Create token with claims
	claims := jwt.MapClaims{
		"user_id": user.ID, // Ensure this is set as float64
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("login", fiber.Map{
			"Title": "Login - Expense Tracker",
			"Error": "Server error occurred",
		})
	}

	cookie := fiber.Cookie{
		Name:     "token",
		Value:    tokenStr,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	}

	c.Cookie(&cookie)
	return c.Redirect("/")
}

func ProcessRegisterForm(c *fiber.Ctx) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")

	// Simple validation
	if name == "" || email == "" || password == "" {
		return c.Render("register", fiber.Map{
			"Title": "Register - Expense Tracker",
			"Error": "All fields are required",
			"Name":  name,
			"Email": email,
		})
	}

	if password != confirmPassword {
		return c.Render("register", fiber.Map{
			"Title": "Register - Expense Tracker",
			"Error": "Passwords do not match",
			"Name":  name,
			"Email": email,
		})
	}

	// Check if email already exists
	var existingUser models.User
	if result := config.DB.Where("email = ?", email).First(&existingUser); result.RowsAffected > 0 {
		return c.Render("register", fiber.Map{
			"Title": "Register - Expense Tracker",
			"Error": "Email already in use",
			"Name":  name,
			"Email": email,
		})
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		return c.Render("register", fiber.Map{
			"Title": "Register - Expense Tracker",
			"Error": "Error creating account. Please try again.",
			"Name":  name,
			"Email": email,
		})
	}

	return c.Redirect("/login?registered=true")
}

// Logout handler
func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	return c.Redirect("/login")
}
