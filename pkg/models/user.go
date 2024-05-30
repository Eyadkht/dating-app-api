package models

import (
	"time"
)

type User struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id" `
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
	Token    Token  `gorm:"constraint:OnDelete:CASCADE;"`
}

type Token struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Value     string    `json:"value"`
	UserID    uint64    `gorm:"unique" json:"userID"`
	CreatedAt time.Time `json:"createdAt"`
}
