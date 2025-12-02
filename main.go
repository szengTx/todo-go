package main

import (
	"html/template"
	"net/http"

	"todo-go/database"
	"todo-go/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	r := gin.Default()
	r.SetHTMLTemplate(template.Must(template.ParseGlob("templates/*.html")))

	// 路由
	r.GET("/register", handlers.Register)
	r.POST("/register", handlers.Register)
	r.GET("/login", handlers.Login)
	r.POST("/login", handlers.Login)
	r.GET("/tasks", handlers.TasksPage)
	r.POST("/tasks/create", handlers.CreateTask)
	r.POST("/tasks/toggle/:id", handlers.ToggleTask)
	r.GET("/tasks/delete/:id", handlers.DeleteTask)

	r.GET("/", func(c *gin.Context) {
		if _, err := c.Cookie("user_id"); err == nil {
			c.Redirect(http.StatusFound, "/tasks")
		} else {
			c.Redirect(http.StatusFound, "/login")
		}
	})

	r.Run(":8080")
}
