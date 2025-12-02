package database

import (
	"todo-go/models" // 注意：这里用你的模块名作为前缀

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移
	DB.AutoMigrate(&models.User{}, &models.Task{})
}
