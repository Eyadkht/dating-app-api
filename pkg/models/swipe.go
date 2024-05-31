package models

import (
	"time"
)

type Swipe struct {
	ID        uint64    `json:"id" gorm:"primary_key"`
	SwiperID  uint64    `json:"swiperID" gorm:"foreignKey:SwiperID;references:UserID"`
	TargetID  uint64    `json:"targetID" gorm:"foreignKey:TargetID;references:UserID"`
	SwipeType string    `json:"swipeType" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null"`
}

type Match struct {
	ID        uint64    `json:"id" gorm:"primary_key"`
	User1ID   uint64    `json:"user1ID" gorm:"foreignKey:User1ID;references:UserID"`
	User2ID   uint64    `json:"user2ID" gorm:"foreignKey:User2ID;references:UserID"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null"`
}
