package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"todo-go/database"
	"todo-go/models"

	"github.com/gin-gonic/gin"
)

// CalendarEvent 用于为 FullCalendar 格式化任务数据
type CalendarEvent struct {
	Title  string `json:"title"`
	Start  string `json:"start"`
	End    string `json:"end,omitempty"`
	AllDay bool   `json:"allDay"`
}

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
	tmpl.ExecuteTemplate(c.Writer, "base", gin.H{
		"Tasks":      tasks,
		"IsLoggedIn": true,
	})
}

func GetTasksJSON(c *gin.Context) {
	cookie, _ := c.Cookie("user_id")
	userID, _ := strconv.Atoi(cookie)

	var tasks []models.Task
	database.DB.Where("user_id = ?", userID).Find(&tasks)

	events := make([]CalendarEvent, 0)
	for _, task := range tasks {
		if task.Deadline != nil {
			start := task.Deadline
			end := task.Deadline
			if !task.AllDay {
				// 默认持续 1 小时，便于在 timeGrid 中占位
				oneHour := task.Deadline.Add(time.Hour)
				end = &oneHour
			}
			startStr := start.Format(time.RFC3339)
			endStr := ""
			if end != nil {
				endStr = end.Format(time.RFC3339)
			}
			if task.AllDay {
				startStr = task.Deadline.Format("2006-01-02")
				endStr = ""
			}
			events = append(events, CalendarEvent{
				Title:  task.Title,
				Start:  startStr,
				End:    endStr,
				AllDay: task.AllDay,
			})
		}
	}

	c.JSON(http.StatusOK, events)
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
	deadlineDate := c.PostForm("deadline_date")
	deadlineTime := c.PostForm("deadline_time")
	allDay := c.PostForm("all_day") != ""

	if title == "" {
		// 如果标题为空，重新查询当前用户的任务并在页面中显示错误
		var tasks []models.Task
		database.DB.Where("user_id = ?", userID).Find(&tasks)
		c.Writer.WriteHeader(http.StatusOK)
		tmpl.ExecuteTemplate(c.Writer, "base", gin.H{"error": "任务标题不能为空", "Tasks": tasks})
		return
	}

	var deadline *time.Time
	if deadlineDate != "" {
		layout := "2006-01-02"
		if !allDay && deadlineTime != "" {
			layout = "2006-01-02 15:04"
			combined := deadlineDate + " " + deadlineTime
			if parsedTime, err := time.Parse(layout, combined); err == nil {
				deadline = &parsedTime
			}
		} else {
			if parsedDate, err := time.Parse(layout, deadlineDate); err == nil {
				deadline = &parsedDate
			}
		}
	}

	task := models.Task{Title: title, UserID: uint(userID), Deadline: deadline, AllDay: allDay}
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
