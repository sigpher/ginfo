package model

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Dept     string `json:"dept"`
	Email    string `gorm:"unique" `
	Phone    string `gorm:"unique"`
	Password string `json:"password"`
}
