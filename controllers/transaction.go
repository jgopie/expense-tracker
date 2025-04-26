package controllers

import (
	"expense-tracker/config"
	"expense-tracker/models"

	"github.com/gofiber/fiber/v2"
)

func GetTransactions(c *fiber.Ctx) error {
	var transactions []models.Transaction
	config.DB.Find(&transactions)
	return c.JSON(transactions)
}

func CreateTransaction(c *fiber.Ctx) error {
	transaction := new(models.Transaction)
	if err := c.BodyParser(transaction); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse json"})
	}
	config.DB.Create(&transaction)
	return c.JSON(transaction)
}

func DeleteTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	config.DB.Delete(&models.Transaction{}, id)
	return c.SendStatus(fiber.StatusNoContent)
}
