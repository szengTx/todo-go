package models

import "time"

// Task 定义了任务模型
type Task struct {
	ID        uint       `gorm:"primaryKey"`
	Title     string     `gorm:"not null"`
	Completed bool       `gorm:"default:false"`
	UserID    uint       `gorm:"not null"`
	Deadline  *time.Time `gorm:"null"`
}
