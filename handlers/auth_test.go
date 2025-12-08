package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"todo-go/database"
	"todo-go/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
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

func TestRegisterAndLogin(t *testing.T) {
	setupTestDB(t)

	// Use a real Gin router to ensure request handling/parsing matches runtime
	r := gin.New()
	r.POST("/register", Register)
	r.POST("/login", Login)

	form := url.Values{}
	form.Set("username", "testuser")
	form.Set("password", "pass")

	// Register via router
	req := httptest.NewRequest("POST", "/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusFound {
		t.Fatalf("register expected redirect, got %d body=%s", w.Code, w.Body.String())
	}

	// Login via router
	req = httptest.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusFound {
		t.Fatalf("login expected redirect, got %d body=%s", w.Code, w.Body.String())
	}

	// Cookie should be set
	if len(w.Header().Values("Set-Cookie")) == 0 {
		t.Fatalf("expected Set-Cookie header, got none")
	}

	// DB should contain the user
	var u models.User
	if database.DB.Where("username = ?", "testuser").First(&u).Error != nil {
		t.Fatalf("user not found in db")
	}
}
