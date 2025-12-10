package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"todo-go/database"
	"todo-go/models"

	"github.com/gin-gonic/gin"
)

// ProfilePage 显示并更新个人信息
func ProfilePage(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/profile.html", "templates/base.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing templates: %v", err)
		return
	}

	cookie, err := c.Cookie("user_id")
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	userID, err := strconv.Atoi(cookie)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.String(http.StatusInternalServerError, "User not found")
		return
	}

	if c.Request.Method == http.MethodPost {
		displayName := c.PostForm("display_name")
		email := c.PostForm("email")
		avatarURL := c.PostForm("avatar_url")

		user.DisplayName = displayName
		user.Email = email
		user.AvatarURL = avatarURL

		if err := database.DB.Save(&user).Error; err != nil {
			c.Writer.WriteHeader(http.StatusOK)
			tmpl.ExecuteTemplate(c.Writer, "base", gin.H{
				"User":       user,
				"IsLoggedIn": true,
				"error":      "更新失败，请稍后重试",
			})
			return
		}

		tmpl.ExecuteTemplate(c.Writer, "base", gin.H{
			"User":       user,
			"IsLoggedIn": true,
			"success":    "已更新",
		})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	tmpl.ExecuteTemplate(c.Writer, "base", gin.H{
		"User":       user,
		"IsLoggedIn": true,
	})
}

// UploadAvatar 处理头像上传，保存到 static/uploads 并更新用户头像 URL
func UploadAvatar(c *gin.Context) {
	cookie, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	userID, err := strconv.Atoi(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未选择文件"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持 png/jpg/jpeg"})
		return
	}

	if err := os.MkdirAll("static/uploads", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建目录"})
		return
	}

	filename := fmt.Sprintf("avatar_%d_%d%s", userID, time.Now().Unix(), ext)
	path := filepath.Join("static", "uploads", filename)

	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err == nil {
		user.AvatarURL = "/static/uploads/" + filename
		database.DB.Save(&user)
	}

	c.JSON(http.StatusOK, gin.H{"url": "/static/uploads/" + filename})
}
