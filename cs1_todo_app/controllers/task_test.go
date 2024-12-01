package controllers

import (
	"bytes"
	"context"
	"cs1_todo_app/database"
	"cs1_todo_app/models"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func setupTestDB() {
	ctx := context.Background()

	err := godotenv.Load("../.env")
	if err != nil {
		logger.Fatal("Error loading .env file: ", zap.Error(err))
	}

	dbname := os.Getenv("DB_DATABASE")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	if dbname == "" || username == "" || password == "" {
		logger.Fatal("Database environment variables are not set properly")
	}

	// Create PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       dbname,
			"POSTGRES_USER":     username,
			"POSTGRES_PASSWORD": password,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		panic(err)
	}

	// Get the container's host and port
	host, _ := postgresC.Host(ctx)
	logger.Info("Host: " + host)
	//port, _ := postgresC.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, username, password, dbname)

	// Connect to the PostgreSQL database
	database.InitDB(dsn)

	// Tear down the container after tests
	defer func(postgresC testcontainers.Container, ctx context.Context) {
		err := postgresC.Terminate(ctx)
		if err != nil {

		}
	}(postgresC, ctx)
}

// Test the GetAllTasks controller
func TestGetAllTasks(t *testing.T) {
	// Set up the test database (make sure setupTestDB initializes the DB properly)
	setupTestDB()

	// Ensure database connection is closed after the test
	defer database.CloseDB()

	_, err := database.GetDB().Exec("DELETE FROM task")
	if err != nil {
		t.Fatalf("Failed to clear tasks: %v", err)
	}

	// Insert two tasks into the database
	_, err = database.GetDB().Exec(`INSERT INTO task(task, completed, created_at, updated_at) 
		VALUES($1, $2, $3, $4)`, "Task 1", false, time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to insert task 1: %v", err)
	}

	_, err = database.GetDB().Exec(`INSERT INTO task(task, completed, created_at, updated_at) 
		VALUES($1, $2, $3, $4)`, "Task 2", false, time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to insert task 2: %v", err)
	}

	// Set up the HTTP request
	req := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()

	// Call the GetAllTasks controller
	GetAllTasks(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Unmarshal the JSON response into a slice of tasks
	var tasks []models.Task
	err = json.Unmarshal(w.Body.Bytes(), &tasks)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check that the response contains two tasks
	assert.Len(t, tasks, 2)

	// Verify the content of the tasks
	assert.Equal(t, "Task 1", tasks[0].Task)
	assert.Equal(t, "Task 2", tasks[1].Task)
}

func TestGetTaskByID(t *testing.T) {
	// Set up the test database (make sure setupTestDB initializes the DB properly)
	setupTestDB()

	// Ensure database connection is closed after the test
	defer database.CloseDB()

	_, err := database.GetDB().Exec("DELETE FROM task")
	if err != nil {
		t.Fatalf("Failed to clear tasks: %v", err)
	}

	// Insert one task into the database
	_, err = database.GetDB().Exec(`INSERT INTO task(task_id, task, completed, created_at, updated_at) 
		VALUES($1, $2, $3, $4, $5)`, 1, "Test Task 1", false, time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to insert task 1: %v", err)
	}

	// Set up the HTTP request
	req := httptest.NewRequest("GET", "/tasks/1", nil)
	w := httptest.NewRecorder()

	// Call the GetAllTasks controller
	GetTaskByID(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Unmarshal the JSON response into a slice of tasks
	var task models.Task
	err = json.Unmarshal(w.Body.Bytes(), &task)
	assert.NoError(t, err)
	assert.Equal(t, "Test Task 1", task.Task)
}

func TestCreateNewTask(t *testing.T) {
	// Set up the test database (make sure setupTestDB initializes the DB properly)
	setupTestDB()

	// Ensure database connection is closed after the test
	defer database.CloseDB()

	_, err := database.GetDB().Exec("DELETE FROM task")
	if err != nil {
		t.Fatalf("Failed to clear tasks: %v", err)
	}

	task := map[string]interface{}{
		"task":      "New Task",
		"completed": false,
	}
	body, _ := json.Marshal(task)
	// Set up the HTTP request
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	CreateNewTask(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)
	var createdTask models.Task
	err = json.Unmarshal(w.Body.Bytes(), &createdTask)
	assert.NoError(t, err)
	assert.Equal(t, "New Task", createdTask.Task)
}

func TestUpdateTask(t *testing.T) {
	// Set up the test database (make sure setupTestDB initializes the DB properly)
	setupTestDB()

	// Ensure database connection is closed after the test
	defer database.CloseDB()

	_, err := database.GetDB().Exec("DELETE FROM task")
	if err != nil {
		t.Fatalf("Failed to clear tasks: %v", err)
	}

	// Insert one task into the database
	_, err = database.GetDB().Exec(`INSERT INTO task(task_id, task, completed, created_at, updated_at) 
		VALUES($1, $2, $3, $4, $5)`, 1, "Test Task 1", false, time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to insert task 1: %v", err)
	}

	updateData := map[string]interface{}{
		"task":      "Updated Task",
		"completed": true,
	}
	body, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	UpdateTask(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	var updatedTask models.Task
	err = json.Unmarshal(w.Body.Bytes(), &updatedTask)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Task", updatedTask.Task)
}

func TestDeleteTask(t *testing.T) {
	// Set up the test database (make sure setupTestDB initializes the DB properly)
	setupTestDB()

	// Ensure database connection is closed after the test
	defer database.CloseDB()

	_, err := database.GetDB().Exec("DELETE FROM task")
	if err != nil {
		t.Fatalf("Failed to clear tasks: %v", err)
	}

	// Insert one task into the database
	_, err = database.GetDB().Exec(`INSERT INTO task(task_id, task, completed, created_at, updated_at) 
		VALUES($1, $2, $3, $4, $5)`, 1, "Test Task 1", false, time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to insert task 1: %v", err)
	}

	// Create test request
	req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
	w := httptest.NewRecorder()

	// Perform request
	DeleteTask(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Task removed", response["message"])
}
