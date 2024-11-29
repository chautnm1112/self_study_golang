package controllers

import (
	"CS1_ToDoApp/database"
	"CS1_ToDoApp/models"
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"
)

func extractIDFromURL(path string) string {
	// Simple regex to match /tasks/<id>
	re := regexp.MustCompile(`/tasks/(\d+)`)
	match := re.FindStringSubmatch(path)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func setupRouter() http.Handler {
	// Create a new ServeMux instance to avoid global route conflicts
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetAllTasks(w, r)
		} else if r.Method == http.MethodPost {
			CreateNewTask(w, r)
		}
	})

	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		id := extractIDFromURL(r.URL.Path)

		if id == "" {
			http.Error(w, "Task ID is required", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			GetTaskByID(w, r)
		} else if r.Method == http.MethodPut {
			UpdateTask(w, r)
		} else if r.Method == http.MethodDelete {
			DeleteTask(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}

func TestGetAllTasks(t *testing.T) {
	// Mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Replace the global DB connection with the mock connection
	database.Db = db

	// Mock rows
	mockRows := sqlmock.NewRows([]string{"task_id", "task", "completed", "created_at", "updated_at"}).
		AddRow(1, "Task 1", false, time.Now(), time.Now()).
		AddRow(2, "Task 2", true, time.Now(), time.Now())
	mock.ExpectQuery("SELECT task_id, task, completed, created_at, updated_at FROM task").WillReturnRows(mockRows)

	// Set up router and routes
	r := setupRouter()
	//r.GET("/tasks", GetAllTasks)

	// Create test request
	req, _ := http.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	var tasks []models.Task
	err = json.Unmarshal(w.Body.Bytes(), &tasks)
	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, "Task 1", tasks[0].Task)
	assert.Equal(t, "Task 2", tasks[1].Task)

	// Ensure all expectations met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetTaskByID(t *testing.T) {
	// Mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Replace the global DB connection with the mock connection
	database.Db = db

	// Mock row
	mockRow := sqlmock.NewRows([]string{"task_id", "task", "completed", "created_at", "updated_at"}).
		AddRow(1, "Test Task", false, time.Now(), time.Now())
	mock.ExpectQuery("SELECT task_id, task, completed, created_at, updated_at FROM task WHERE task_id = \\$1").
		WithArgs("1").
		WillReturnRows(mockRow)

	// Set up router and routes
	r := setupRouter()
	//r.GET("/tasks/:id", GetTaskByID)

	// Create test request
	req, _ := http.NewRequest("GET", "/tasks/1", nil)
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	var task models.Task
	err = json.Unmarshal(w.Body.Bytes(), &task)
	assert.NoError(t, err)
	assert.Equal(t, "Test Task", task.Task)

	// Ensure all expectations met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCreateNewTask(t *testing.T) {
	// Mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Replace the global DB connection with the mock connection
	database.Db = db

	// Mock query
	mock.ExpectQuery("INSERT INTO task\\(task, completed, created_at, updated_at\\) VALUES\\(\\$1, \\$2, \\$3, \\$4\\) RETURNING task_id").
		WithArgs("New Task", false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"task_id"}).AddRow(1))

	// Set up router and routes
	r := setupRouter()
	//r.POST("/tasks", CreateNewTask)

	// Create test request
	task := map[string]interface{}{
		"task":      "New Task",
		"completed": false,
	}
	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)
	var createdTask models.Task
	err = json.Unmarshal(w.Body.Bytes(), &createdTask)
	assert.NoError(t, err)
	assert.Equal(t, "New Task", createdTask.Task)

	// Ensure all expectations met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateTask(t *testing.T) {
	// Mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Replace the global DB connection with the mock connection
	database.Db = db

	// Mock initial query to fetch task
	mock.ExpectQuery("SELECT task_id, task, completed, created_at, updated_at FROM task WHERE task_id = \\$1").
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"task_id", "task", "completed", "created_at", "updated_at"}).
			AddRow(1, "Old Task", false, time.Now(), time.Now()))

	// Mock update query
	mock.ExpectExec("UPDATE task SET task = \\$1, completed = \\$2, updated_at = \\$3 WHERE task_id = \\$4").
		WithArgs("Updated Task", true, sqlmock.AnyArg(), "1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Set up router and routes
	r := setupRouter()
	//r.PUT("/tasks/:id", UpdateTask)

	// Create test request
	updateData := map[string]interface{}{
		"task":      "Updated Task",
		"completed": true,
	}
	body, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	var updatedTask models.Task
	err = json.Unmarshal(w.Body.Bytes(), &updatedTask)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Task", updatedTask.Task)

	// Ensure all expectations met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeleteTask(t *testing.T) {
	// Mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Replace the global DB connection with the mock connection
	database.Db = db

	// Mock delete query
	mock.ExpectExec("DELETE FROM task WHERE task_id = \\$1").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Set up router and routes
	r := setupRouter()
	//r.DELETE("/tasks/:id", DeleteTask)

	// Create test request
	req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Task removed", response["message"])

	// Ensure all expectations met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
