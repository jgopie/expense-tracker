package controllers

import (
	"encoding/csv"
	"expense-tracker/config"
	"expense-tracker/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ExportTransactions(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(float64)
	var transactions []models.Transaction
	if err := config.DB.Preload("Account").
		Where("user_id = ?", uint(userId)).
		Find(&transactions).Error; err != nil {
		return c.Status(500).SendString("error fetching transactions")
	}
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=transactions.csv")
	writer := csv.NewWriter(c.Response().BodyWriter())
	defer writer.Flush()

	header := []string{
		"Date",
		"Account",
		"Description",
		"Amount",
		"Category",
	}
	if err := writer.Write(header); err != nil {
		return err
	}
	for _, t := range transactions {
		record := []string{
			t.CreatedAt.Format("2006-01-02"),
			t.Account.Name,
			t.Description,
			fmt.Sprintf("%.2f", t.Amount),
			t.Category,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}
