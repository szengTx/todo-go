package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"todo-go/database"
	"todo-go/models"

	"github.com/gin-gonic/gin"
)

func TasksPage(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/tasks.html", "templates/base.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error parsing templates: %v", err)
		return
	}
	cookie, _ := c.Cookie("user_id")
	userID, _ := strconv.Atoi(cookie)

	var tasks []models.Task
	database.DB.Where("user_id = ?", userID).Find(&tasks)

	c.Writer.WriteHeader(http.StatusOK)
	tmpl.ExecuteTemplate(c.Writer, "base", gin.H{"Tasks": tasks})
}

func CreateTask(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/tasks.html", "templates/base.html")
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

	title := c.PostForm("title")
	if title == "" {
		// 如果标题为空，重新查询当前用户的任务并在页面中显示错误
		var tasks []models.Task
		database.DB.Where("user_id = ?", userID).Find(&tasks)
		c.Writer.WriteHeader(http.StatusOK)
		tmpl.ExecuteTemplate(c.Writer, "base", gin.H{"error": "任务标题不能为空", "Tasks": tasks})
		return
	}

	task := models.Task{Title: title, UserID: uint(userID)}
	if err := database.DB.Create(&task).Error; err != nil {
		var tasks []models.Task
		database.DB.Where("user_id = ?", userID).Find(&tasks)
		c.Writer.WriteHeader(http.StatusOK)
		tmpl.ExecuteTemplate(c.Writer, "base", gin.H{"error": "无法创建任务", "Tasks": tasks})
		return
	}
	c.Redirect(http.StatusFound, "/tasks")
}

// ToggleTask 切换任务的完成状态
func ToggleTask(c *gin.Context) {
	id := c.Param("id")
	cookie, _ := c.Cookie("user_id")
	userID, _ := strconv.Atoi(cookie)

	var task models.Task
	if database.DB.Where("id = ? AND user_id = ?", id, userID).First(&task).Error == nil {
		task.Completed = !task.Completed
		database.DB.Save(&task)
	}
	c.Redirect(http.StatusFound, "/tasks")
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	cookie, _ := c.Cookie("user_id")
	userID, _ := strconv.Atoi(cookie)

	database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Task{})
	c.Redirect(http.StatusFound, "/tasks")
}
