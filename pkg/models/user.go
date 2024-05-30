package models

type User struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id" `
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}
