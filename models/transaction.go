package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UserId      uint      `json:"user_id"` // Foreign Key, links to user table
	AccountId   uint      `json:"account_id"`
	Account     Account   `json:"account" gorm:"foreignKey:AccountId"`
}
