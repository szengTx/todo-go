package models

type User struct {
	ID          uint   `gorm:"primaryKey"`
	Username    string `gorm:"unique;not null"`
	Password    string `gorm:"not null"`
	DisplayName string `gorm:"size:100"`
	Email       string `gorm:"size:120"`
	AvatarURL   string `gorm:"size:255"`
}
