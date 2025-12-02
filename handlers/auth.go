package handlers

import (
	"fmt"
	"net/http"

	"todo-go/database"
	"todo-go/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "register.html", nil)
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")

	var existingUser models.User
	if database.DB.Where("username = ?", username).First(&existingUser).Error == nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "用户名已存在"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "注册失败"})
		return
	}

	user := models.User{Username: username, Password: string(hashedPassword)}
	database.DB.Create(&user)

	c.Redirect(http.StatusFound, "/login")
}

func Login(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "login.html", nil)
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")

	var user models.User
	if database.DB.Where("username = ?", username).First(&user).Error != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "用户名或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "用户名或密码错误"})
		return
	}

	c.SetCookie("user_id", fmt.Sprintf("%d", user.ID), 3600, "/", "", false, true)
	c.Redirect(http.StatusFound, "/tasks")
}
