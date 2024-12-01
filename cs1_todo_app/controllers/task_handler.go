package controllers

import (
	"cs1_todo_app/database"
	"cs1_todo_app/models"
	"database/sql"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Get the global logger instance
var logger, _ = zap.NewProduction()

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	query := "SELECT task_id, task, completed, created_at, updated_at FROM task"

	// Log truy vấn
	logger.Info("Executing query", zap.String("query", query))

	rows, err := database.GetDB().Query(query)

	if err != nil {
		logger.Error("Query execution failed", zap.String("query", query), zap.Error(err))
		err := json.NewEncoder(w).Encode(models.ApiResponse[[]models.Task]{Code: http.StatusInternalServerError, Data: nil, Message: "Database query failed"})
		if err != nil {
			return
		}
		return
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Task, &task.Completed, &task.CreatedAt, &task.UpdatedAt); err != nil {
			logger.Error("Error scanning row", zap.Error(err))
			err := json.NewEncoder(w).Encode(models.ApiResponse[[]models.Task]{Code: http.StatusInternalServerError, Data: nil, Message: "Error when scanning tasks"})
			if err != nil {
				return
			}
			return
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Error iterating over rows", zap.Error(err))
		err := json.NewEncoder(w).Encode(models.ApiResponse[[]models.Task]{Code: http.StatusInternalServerError, Data: nil, Message: "Error when processing tasks"})
		if err != nil {
			return
		}
		return
	}

	err = json.NewEncoder(w).Encode(models.ApiResponse[[]models.Task]{Code: http.StatusOK, Data: &tasks, Message: "Get all task success"})
	if err != nil {
		return
	}
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/tasks/"):]

	var task models.Task
	row := database.GetDB().QueryRow("SELECT task_id, task, completed, created_at, updated_at FROM task WHERE task_id = $1", id)
	if err := row.Scan(&task.ID, &task.Task, &task.Completed, &task.CreatedAt, &task.UpdatedAt); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			err := json.NewEncoder(w).Encode(models.ApiResponse[[]models.Task]{Code: http.StatusNotFound, Data: nil, Message: "Task not found"})
			if err != nil {
				return
			}
		} else {
			logger.Error("Error querying task", zap.Error(err))
			err := json.NewEncoder(w).Encode(models.ApiResponse[[]models.Task]{Code: http.StatusInternalServerError, Data: nil, Message: "Error querying task"})
			if err != nil {
				return
			}
		}
		return
	}

	err := json.NewEncoder(w).Encode(models.ApiResponse[models.Task]{Code: http.StatusOK, Data: &task, Message: "Get task success"})
	if err != nil {
		return
	}
}

func CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		// Log error with zap
		logger.Error("Invalid JSON format", zap.Error(err))
		err := json.NewEncoder(w).Encode(models.ApiResponse[[]models.Task]{Code: http.StatusBadRequest, Data: nil, Message: "Invalid JSON format"})
		if err != nil {
			return
		}
		return
	}

	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	var lastInsertID int64
	err := database.GetDB().QueryRow(`INSERT INTO task(task, completed, created_at, updated_at)
		VALUES($1, $2, $3, $4) RETURNING task_id`, task.Task, task.Completed, task.CreatedAt, task.UpdatedAt).Scan(&lastInsertID)
	if err != nil {
		// Log error with zap
		logger.Error("Can't create task", zap.Error(err))
		err := json.NewEncoder(w).Encode(models.ApiResponse[[]models.Task]{Code: http.StatusInternalServerError, Data: nil, Message: "Can't create task"})
		if err != nil {
			return
		}
		return
	}

	task.ID = lastInsertID
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(models.ApiResponse[models.Task]{Code: http.StatusCreated, Data: &task, Message: "Create new task success"})
	if err != nil {
		return
	}
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/tasks/"):]

	var task models.Task
	if err := database.GetDB().QueryRow("SELECT task_id, task, completed, created_at, updated_at FROM task WHERE task_id = $1", id).Scan(&task.ID, &task.Task, &task.Completed, &task.CreatedAt, &task.UpdatedAt); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			err := json.NewEncoder(w).Encode(models.ApiResponse[[]models.Task]{Code: http.StatusNotFound, Data: nil, Message: "Task not found"})
			if err != nil {
				return
			}
		} else {
			logger.Error("Error querying task", zap.Error(err))
			err := json.NewEncoder(w).Encode(models.ApiResponse[[]models.Task]{Code: http.StatusInternalServerError, Data: nil, Message: "Error querying task"})
			if err != nil {
				return
			}
		}
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		logger.Error("Invalid JSON format", zap.Error(err))
		err := json.NewEncoder(w).Encode(models.ApiResponse[[]models.Task]{Code: http.StatusBadRequest, Data: nil, Message: "Invalid JSON format"})
		if err != nil {
			return
		}
		return
	}

	task.UpdatedAt = time.Now()

	_, err := database.GetDB().Exec(`UPDATE task SET task = $1, completed = $2, updated_at = $3 WHERE task_id = $4`,
		task.Task, task.Completed, task.UpdatedAt, id)
	if err != nil {
		logger.Error("Can't update task", zap.Error(err))
		err := json.NewEncoder(w).Encode(models.ApiResponse[[]models.Task]{Code: http.StatusInternalServerError, Data: nil, Message: "Can't update task"})
		if err != nil {
			return
		}
		return
	}

	err = json.NewEncoder(w).Encode(models.ApiResponse[models.Task]{Code: http.StatusOK, Data: &task, Message: "Update the task success"})
	if err != nil {
		return
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/tasks/"):]

	result, err := database.GetDB().Exec("DELETE FROM task WHERE task_id = $1", id)
	if err != nil {
		// Log error with zap
		logger.Error("Can't remove task", zap.Error(err))
		err := json.NewEncoder(w).Encode(models.ApiResponse[[]models.Task]{Code: http.StatusInternalServerError, Data: nil, Message: "Can't remove task"})
		if err != nil {
			return
		}
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		err := json.NewEncoder(w).Encode(models.ApiResponse[[]models.Task]{Code: http.StatusNotFound, Data: nil, Message: "Task not found"})
		if err != nil {
			return
		}
		return
	}

	logger.Info("Task removed", zap.String("task_id", id))

	err = json.NewEncoder(w).Encode(models.ApiResponse[models.Task]{Code: http.StatusCreated, Data: nil, Message: "Task removed"})
	if err != nil {
		return
	}
}
