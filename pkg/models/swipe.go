package models

import (
	"time"
)

type Swipe struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	SwiperID  uint      `json:"swiperID" gorm:"foreignKey:SwiperID;references:UserID"`
	TargetID  uint      `json:"targetID" gorm:"foreignKey:TargetID;references:UserID"`
	SwipeType string    `json:"swipeType" gorm:"type:enum('like','dislike');not null"`
	SwipedAt  time.Time `json:"swipedAt" gorm:"not null"`
}
type Match struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	User1ID   uint      `json:"user1ID" gorm:"foreignKey:User1ID;references:UserID"`
	User2ID   uint      `json:"user2ID" gorm:"foreignKey:User2ID;references:UserID"`
	MatchedAt time.Time `json:"matchedAt" gorm:"not null"`
}
