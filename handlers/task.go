package handlers

import (
	"net/http"
	"strconv"

	"todo-go/database"
	"todo-go/models"

	"github.com/gin-gonic/gin"
)

func TasksPage(c *gin.Context) {
	cookie, _ := c.Cookie("user_id")
	userID, _ := strconv.Atoi(cookie)

	var tasks []models.Task
	database.DB.Where("user_id = ?", userID).Find(&tasks)

	c.HTML(http.StatusOK, "tasks.html", gin.H{"Tasks": tasks})
}

func CreateTask(c *gin.Context) {
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
		c.HTML(http.StatusOK, "tasks.html", gin.H{"error": "任务标题不能为空"})
		return
	}

	task := models.Task{Title: title, UserID: uint(userID)}
	database.DB.Create(&task)
	c.Redirect(http.StatusFound, "/tasks")
}

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
