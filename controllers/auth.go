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

func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: string(password),
	}
	result := config.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "error with user registration"})
	}
	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	var user models.User
	config.DB.Where("email = ?", data["email"]).First(&user)

	if user.ID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid username or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid username or password"})
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
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

	return c.JSON(fiber.Map{"token": tokenStr})
}

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
	config.DB.Where("email = ?", email).First(&user)
	if user.ID == 0 {
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

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
	return c.Redirect("/dashboard")
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
