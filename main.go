package main

import (
	"net/http"

	"todo-go/database"
	"todo-go/handlers"
	"todo-go/models"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	r := gin.Default()
	r.Static("/static", "./static")
	r.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Next()
	})
	// 根路由：如果没有任何用户，提示先注册，否则去登录
	r.GET("/", func(c *gin.Context) {
		var u models.User
		if database.DB.First(&u).Error != nil {
			// 数据库中没有用户，去注册
			c.Redirect(http.StatusFound, "/register")
			return
		}
		c.Redirect(http.StatusFound, "/login")
	})

	// 公开路由
	r.GET("/register", handlers.Register)
	r.POST("/register", handlers.Register)
	r.GET("/login", handlers.Login)
	r.POST("/login", handlers.Login)
	r.GET("/logout", handlers.Logout)

	// 受保护路由
	auth := r.Group("/tasks")
	auth.Use(handlers.AuthMiddleware())
	{
		auth.GET("/", handlers.TasksPage)
		auth.GET("/json", handlers.GetTasksJSON)
		auth.POST("/create", handlers.CreateTask)
		auth.POST("/toggle/:id", handlers.ToggleTask)
		auth.GET("/delete/:id", handlers.DeleteTask)
	}

	profile := r.Group("/")
	profile.Use(handlers.AuthMiddleware())
	{
		profile.GET("/profile", handlers.ProfilePage)
		profile.POST("/profile", handlers.ProfilePage)
		profile.POST("/profile/upload", handlers.UploadAvatar) // 上传头像
	}

	r.Run(":8080")
}
