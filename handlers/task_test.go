package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"todo-go/database"
	"todo-go/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBForTasks(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	database.DB = db
	if err := database.DB.AutoMigrate(&models.User{}, &models.Task{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	gin.SetMode(gin.TestMode)
}

func TestCreateToggleDeleteTask(t *testing.T) {
	setupTestDBForTasks(t)

	// create a user directly
	u := models.User{Username: "u1", Password: "x"}
	if err := database.DB.Create(&u).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	userID := int(u.ID)

	// Create task via router to ensure request parsing
	r := gin.New()
	r.POST("/tasks/create", CreateTask)
	r.POST("/tasks/toggle/:id", ToggleTask)
	r.GET("/tasks/delete/:id", DeleteTask)

	form := url.Values{}
	form.Set("title", "Task1")
	req, _ := http.NewRequest("POST", "/tasks/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "user_id", Value: strconv.Itoa(userID)})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusFound {
		t.Fatalf("create expected redirect, got %d body=%s", w.Code, w.Body.String())
	}

	var tasks []models.Task
	database.DB.Where("user_id = ?", userID).Find(&tasks)
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	task := tasks[0]

	// Toggle
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/tasks/toggle/"+strconv.Itoa(int(task.ID)), nil)
	req.AddCookie(&http.Cookie{Name: "user_id", Value: strconv.Itoa(userID)})
	r.ServeHTTP(w, req)
	var tsk models.Task
	database.DB.First(&tsk, task.ID)
	if !tsk.Completed {
		t.Fatalf("expected task completed true after toggle")
	}

	// Delete
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/tasks/delete/"+strconv.Itoa(int(task.ID)), nil)
	req.AddCookie(&http.Cookie{Name: "user_id", Value: strconv.Itoa(userID)})
	r.ServeHTTP(w, req)

	var count int64
	database.DB.Model(&models.Task{}).Where("user_id = ?", userID).Count(&count)
	if count != 0 {
		t.Fatalf("expected 0 tasks after delete, got %d", count)
	}
}
