package config

import (
	"expense-tracker/models"
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

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db
	db.AutoMigrate(&models.Transaction{}, &models.User{})
}
