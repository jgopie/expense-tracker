package controllers

import (
	"expense-tracker/config"
	"expense-tracker/models"

	"github.com/gofiber/fiber/v2"
)

func GetTransactions(c *fiber.Ctx) error {
	var transactions []models.Transaction
	userID := c.Locals("user_id").(float64)
	config.DB.Where("user_id = ?", uint(userID)).Find(&transactions)
	return c.JSON(transactions)
}

func CreateTransaction(c *fiber.Ctx) error {
	transaction := new(models.Transaction)
	if err := c.BodyParser(transaction); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse json"})
	}
	userID := c.Locals("user_id").(float64)
	transaction.UserId = uint(userID)

	config.DB.Create(&transaction)
	return c.JSON(transaction)
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
