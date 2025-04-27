package controllers

import (
	"expense-tracker/config"
	"expense-tracker/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetTransactions(c *fiber.Ctx) error {
	var transactions []models.Transaction
	userID := c.Locals("user_id").(float64)
	config.DB.Where("user_id = ?", uint(userID)).Find(&transactions)
	return c.JSON(transactions)
}

func CreateTransaction(c *fiber.Ctx) error {
	amount, err := strconv.ParseFloat(c.FormValue("amount"), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).Render("add", fiber.Map{
			"Title": "Add Expense",
			"Error": "Invalid amount",
		})
	}

	userID, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Redirect("/login")
	}

	transaction := models.Transaction{
		Description: c.FormValue("description"),
		Amount:      amount,
		Category:    c.FormValue("category"),
		CreatedAt:   time.Now(), // Or parse date from form
		UserId:      uint(userID),
	}

	if err := config.DB.Create(&transaction).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("add", fiber.Map{
			"Title": "Add Expense",
			"Error": "Failed to save transaction",
		})
	}

	return c.Redirect("/")
}

func DeleteTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(float64)

	var transaction models.Transaction
	config.DB.First(&transaction, id)

	if transaction.UserId != uint(userID) {
		c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "unauthorized"})
	}

	config.DB.Delete(&transaction)
	return c.SendStatus(fiber.StatusNoContent)
}

func RenderAddTransactionPage(c *fiber.Ctx) error {
	return c.Render("add", fiber.Map{
		"Title": "Add Transaction - Expense Tracker",
	})
}
