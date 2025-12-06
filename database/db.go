// database/db.go
package database

import (
	"todo-go/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("file:todo.db?cache=shared&_fk=1"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	DB.AutoMigrate(&models.User{}, &models.Task{})
}
