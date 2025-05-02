package config

import (
	"expense-tracker/models"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/transactions.db"
	}

	// Verify parent directory
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0775); err != nil { // Note: 0775
		log.Fatalf("Failed to create directory '%s': %v", dbDir, err)
	}

	// Try creating empty file first
	if _, err := os.OpenFile(dbPath, os.O_RDONLY|os.O_CREATE, 0664); err != nil {
		log.Fatalf("Failed to create database file: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db
	db.AutoMigrate(&models.Transaction{}, &models.User{}, &models.Account{})
}

func RunMigrations(db *gorm.DB) error {
	// Check if accounts table exists
	var table string
	db.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name='accounts'").Scan(&table)

	if table == "" {
		// Run initial migration
		if err := db.Exec(`ALTER TABLE transactions ADD COLUMN account_id INTEGER`).Error; err != nil {
			return fmt.Errorf("failed to add account_id column: %v", err)
		}

		if err := db.AutoMigrate(&models.Account{}); err != nil {
			return fmt.Errorf("failed to create accounts table: %v", err)
		}

		// Create default account for existing users
		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			return fmt.Errorf("failed to find users: %v", err)
		}

		for _, user := range users {
			defaultAccount := models.Account{
				Name:    "Default Account",
				Balance: 0,
				UserId:  user.ID,
			}
			if err := db.Create(&defaultAccount).Error; err != nil {
				return fmt.Errorf("failed to create default account: %v", err)
			}

			// Associate existing transactions with default account
			if err := db.Model(&models.Transaction{}).
				Where("user_id = ?", user.ID).
				Update("account_id", defaultAccount.ID).Error; err != nil {
				return fmt.Errorf("failed to update transactions: %v", err)
			}
		}
	}
	return nil
}
