package models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Name        string        `json:"name"`
	Balance     float64       `json:"balance"`
	UserId      uint          `json:"user_id"`
	Transaction []Transaction `json:"transactions" gorm:"foreignKey:AccountId;constraint:OnDelete:CASCADE;"`
}
