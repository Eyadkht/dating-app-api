package models

import (
	"time"
)

type Match struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	User1ID   uint      `json:"user1ID" gorm:"foreignKey:User1ID;references:UserID"`
	User2ID   uint      `json:"user2ID" gorm:"foreignKey:User2ID;references:UserID"`
	MatchedAt time.Time `json:"matchedAt" gorm:"not null"`
}
