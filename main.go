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
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	r.SetHTMLTemplate(tmpl)

	// 首页重定向到登录
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})

	// 公开路由
	r.GET("/register", handlers.Register)
	r.POST("/register", handlers.Register)
	r.GET("/login", handlers.Login)
	r.POST("/login", handlers.Login)
	r.GET("/logout", handlers.Logout)

	// 受保护路由
	auth := r.Group("/")
	auth.Use(handlers.AuthMiddleware())
	{
		auth.GET("/tasks", handlers.TasksPage)
		auth.POST("/tasks/create", handlers.CreateTask)
		auth.POST("/tasks/toggle/:id", handlers.ToggleTask)
		auth.GET("/tasks/delete/:id", handlers.DeleteTask)
	}

	r.Run(":8080")
}
