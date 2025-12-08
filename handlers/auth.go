package handlers

import (
	"net/http"
	"strconv"

	"todo-go/database"
	"todo-go/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	if c.Request.Method == "POST" {
		username := c.PostForm("username")
		password := c.PostForm("password")

		if username == "" || password == "" {
			c.HTML(http.StatusOK, "register.html", gin.H{"error": "用户名和密码不能为空"})
			return
		}

		var existingUser models.User
		if database.DB.Where("username = ?", username).First(&existingUser).Error == nil {
			c.HTML(http.StatusOK, "register.html", gin.H{"error": "用户名已存在"})
			return
		}

		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user := models.User{Username: username, Password: string(hashed)}
		database.DB.Create(&user)

		// 注册后自动登录：设置 user_id cookie 并跳转到任务页
		userIDStr := strconv.Itoa(int(user.ID))
		c.SetCookie("user_id", userIDStr, 3600, "/", "", false, true)
		c.Redirect(http.StatusFound, "/tasks")
		return
	}
	c.HTML(http.StatusOK, "register.html", nil)
}

func Login(c *gin.Context) {
	if c.Request.Method == "POST" {
		username := c.PostForm("username")
		password := c.PostForm("password")

		var user models.User
		if database.DB.Where("username = ?", username).First(&user).Error != nil {
			c.HTML(http.StatusOK, "login.html", gin.H{"error": "用户名或密码错误"})
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
			c.HTML(http.StatusOK, "login.html", gin.H{"error": "用户名或密码错误"})
			return
		}

		userIDStr := strconv.Itoa(int(user.ID))
		c.SetCookie("user_id", userIDStr, 3600, "/", "", false, true)
		c.Redirect(http.StatusFound, "/tasks")
		return
	}
	c.HTML(http.StatusOK, "login.html", nil)
}

func Logout(c *gin.Context) {
	c.SetCookie("user_id", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("user_id")
		if err != nil || cookie == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		if _, err := strconv.Atoi(cookie); err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
