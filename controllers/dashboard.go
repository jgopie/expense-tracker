package controllers

import (
	"expense-tracker/config"
	"expense-tracker/models"

	"github.com/gofiber/fiber/v2"
)

func Dashboard(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(float64)

	// Get accounts with balances
	var accounts []models.Account
	if err := config.DB.
		Where("user_id = ?", uint(userID)).
		Find(&accounts).Error; err != nil {
		return c.Status(500).SendString("Error fetching accounts")
	}

	// Get recent transactions with account info
	var transactions []models.Transaction
	if err := config.DB.
		Joins("Account"). // Changed from Preload to Joins for SQLite
		Where("transactions.user_id = ?", uint(userID)).
		Order("transactions.created_at desc").
		Limit(10).
		Find(&transactions).Error; err != nil {
		return c.Status(500).SendString("Error fetching transactions")
	}

	// Calculate totals
	var totalExpenses float64
	categoryTotals := make(map[string]float64)
	for _, t := range transactions {
		totalExpenses += t.Amount
		categoryTotals[t.Category] += t.Amount
	}

	return c.Render("dashboard", fiber.Map{
		"Title":         "Dashboard",
		"Accounts":      accounts,
		"Transactions":  transactions,
		"TotalExpenses": totalExpenses,
		"CategoryData":  categoryTotals,
	})
}
