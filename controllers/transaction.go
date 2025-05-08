package controllers

import (
	"expense-tracker/config"
	"expense-tracker/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetTransactions(c *fiber.Ctx) error {
	var transactions []models.Transaction
	userID := c.Locals("user_id").(float64)
	config.DB.Where("user_id = ?", uint(userID)).Find(&transactions)
	return c.JSON(transactions)
}

// controllers/transaction.go
func CreateTransaction(c *fiber.Ctx) error {
	// Parse form data
	amount, err := strconv.ParseFloat(c.FormValue("amount"), 64)
	if err != nil {
		return c.Status(400).Render("add", fiber.Map{
			"Error":    "Invalid amount format",
			"Accounts": c.Locals("accounts"),
		})
	}

	accountID, err := strconv.ParseUint(c.FormValue("account_id"), 10, 32)
	if err != nil {
		return c.Status(400).Render("add", fiber.Map{
			"Error":    "Invalid account",
			"Accounts": c.Locals("accounts"),
		})
	}

	userID := c.Locals("user_id").(float64)

	// Start database transaction
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		// Create the transaction
		transaction := models.Transaction{
			Description: c.FormValue("description"),
			Amount:      amount,
			Category:    c.FormValue("category"),
			UserId:      uint(userID),
			AccountId:   uint(accountID),
			CreatedAt:   time.Now(),
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// Update account balance
		if err := tx.Model(&models.Account{}).
			Where("id = ?", accountID).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(500).Render("add", fiber.Map{
			"Error":    "Failed to save transaction",
			"Accounts": c.Locals("accounts"),
		})
	}

	return c.Redirect("/")
}

func DeleteTransaction(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(float64)
	transactionID := c.Params("id")

	// Start database transaction
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// First get the transaction with account info
		var transaction models.Transaction
		if err := tx.
			Where("id = ? AND user_id = ?", transactionID, uint(userID)).
			First(&transaction).Error; err != nil {
			return err // Transaction not found or doesn't belong to user
		}

		// Update the account balance
		if err := tx.Model(&models.Account{}).
			Where("id = ?", transaction.AccountId).
			Update("balance", gorm.Expr("balance - ?", transaction.Amount)).Error; err != nil {
			return err
		}

		// Delete the transaction
		if err := tx.Delete(&transaction).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete transaction",
		})
	}

	return c.Redirect("/")
}

func RenderAddTransactionPage(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(float64)

	var accounts []models.Account
	if err := config.DB.Where("user_id = ?", uint(userID)).Find(&accounts).Error; err != nil {
		return c.Status(500).SendString("Error fetching accounts")
	}

	return c.Render("add", fiber.Map{
		"Title":    "Add Transaction",
		"Accounts": accounts,
	})
}
