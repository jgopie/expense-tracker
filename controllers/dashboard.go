package controllers

import (
	"expense-tracker/config"
	"expense-tracker/models"

	"github.com/gofiber/fiber/v2"
)

func Dashboard(c *fiber.Ctx) error {
	// Get user_id from locals
	userID, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Redirect("/login")
	}

	// Get transactions for this user
	var transactions []models.Transaction
	if err := config.DB.Where("user_id = ?", uint(userID)).Order("created_at desc").Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Server error")
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
		"Transactions":  transactions,
		"TotalExpenses": totalExpenses,
		"CategoryData":  categoryTotals,
	})
}
