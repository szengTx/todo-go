package models

type Task struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Completed bool   `gorm:"default:false"`
	UserID    uint   // 外键，关联 User.ID
}
