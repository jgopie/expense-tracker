package controllers

import (
	"expense-tracker/config"
	"expense-tracker/models"

	"github.com/gofiber/fiber/v2"
)

func GetAccounts(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(float64)

	var accounts []models.Account
	if err := config.DB.Where("user_id = ?", uint(userID)).Find(&accounts).Error; err != nil {
		return c.Status(500).SendString("Error fetching accounts")
	}

	// Calculate total balance
	var totalBalance float64
	for _, account := range accounts {
		totalBalance += account.Balance
	}

	return c.Render("accounts", fiber.Map{
		"Title":        "Accounts",
		"Accounts":     accounts,
		"TotalBalance": totalBalance,
	})
}

func CreateAccount(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(float64)

	account := new(models.Account)
	if err := c.BodyParser(account); err != nil {
		return c.Status(400).SendString("Bad Request")
	}

	account.UserId = uint(userID)
	if err := config.DB.Create(&account).Error; err != nil {
		return c.Status(500).SendString("Error creating account")
	}

	return c.Redirect("/accounts")
}

func DeleteAccount(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(float64)
	accountId := c.Params("id")
	result := config.DB.Where("id = ? AND user_id = ?", accountId, uint(userId)).Delete(&models.Account{})
	if result.Error != nil {
		return c.Status(500).SendString("error deleting account")
	}
	if result.RowsAffected == 0 {
		return c.Status(404).SendString("account not found")
	}
	return c.Redirect("/accounts")
}
