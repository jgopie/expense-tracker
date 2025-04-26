package config

import (
	"expense-tracker/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	db, err := gorm.Open(sqlite.Open("transactions.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to db")
	}
	DB = db
	db.AutoMigrate(&models.Transaction{}, &models.User{})
}
