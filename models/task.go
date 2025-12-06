package models

type Task struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Completed bool
	UserID    uint `gorm:"not null"`
}
