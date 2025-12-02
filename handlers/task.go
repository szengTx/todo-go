package handlers

import (
	"net/http"
	"strconv"

	"todo-go/database"
	"todo-go/models"

	"github.com/gin-gonic/gin"
)

func getCurrentUserID(c *gin.Context) (uint, bool) {
	cookie, err := c.Cookie("user_id")
	if err != nil {
		return 0, false
	}
	id, err := strconv.ParseUint(cookie, 10, 32)
	if err != nil {
		return 0, false
	}
	return uint(id), true
}

func TasksPage(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var tasks []models.Task
	database.DB.Where("user_id = ?", userID).Find(&tasks)

	c.HTML(http.StatusOK, "tasks.html", gin.H{"Tasks": tasks})
}

func CreateTask(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	title := c.PostForm("title")
	if title == "" {
		c.Redirect(http.StatusFound, "/tasks")
		return
	}

	task := models.Task{Title: title, UserID: userID}
	database.DB.Create(&task)
	c.Redirect(http.StatusFound, "/tasks")
}

func ToggleTask(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	taskID := c.Param("id")
	var task models.Task
	if database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	task.Completed = !task.Completed
	database.DB.Save(&task)

	if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
		c.JSON(http.StatusOK, gin.H{"completed": task.Completed})
	} else {
		c.Redirect(http.StatusFound, "/tasks")
	}
}

func DeleteTask(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	taskID := c.Param("id")
	database.DB.Where("id = ? AND user_id = ?", taskID, userID).Delete(&models.Task{})
	c.Redirect(http.StatusFound, "/tasks")
}
