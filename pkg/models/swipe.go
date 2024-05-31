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
