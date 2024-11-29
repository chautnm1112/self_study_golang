package controllers

import (
	"CS1_ToDoApp/database"
	"CS1_ToDoApp/models"
	"database/sql"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	rows, err := database.Db.Query("SELECT task_id, task, completed, created_at, updated_at FROM task")
	if err != nil {
		logrus.Error("Cann't get all task: ", err)
		http.Error(w, "Can't get all tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Task, &task.Completed, &task.CreatedAt, &task.UpdatedAt); err != nil {
			http.Error(w, "Error when getting all tasks", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	json.NewEncoder(w).Encode(tasks)

}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/tasks/"):]
	var task models.Task
	row := database.Db.QueryRow("SELECT task_id, task, completed, created_at, updated_at FROM task WHERE task_id = $1", id)
	if err := row.Scan(&task.ID, &task.Task, &task.Completed, &task.CreatedAt, &task.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error querying task", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(task)
}

func CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	var lastInsertID int64
	err := database.Db.QueryRow(`INSERT INTO task(task, completed, created_at, updated_at)
		VALUES($1, $2, $3, $4) RETURNING task_id`, task.Task, task.Completed, task.CreatedAt, task.UpdatedAt).Scan(&lastInsertID)
	if err != nil {
		http.Error(w, "Cann't create task", http.StatusInternalServerError)
		return
	}

	task.ID = lastInsertID
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)

}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/tasks/"):]
	var task models.Task
	if err := database.Db.QueryRow("SELECT task_id, task, completed, created_at, updated_at FROM task WHERE task_id = $1", id).Scan(&task.ID, &task.Task, &task.Completed, &task.CreatedAt, &task.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error querying task", http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	task.UpdatedAt = time.Now()

	_, err := database.Db.Exec(`UPDATE task SET task = $1, completed = $2, updated_at = $3 WHERE task_id = $4`,
		task.Task, task.Completed, task.UpdatedAt, id)
	if err != nil {
		http.Error(w, "Cann't update task", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/tasks/"):]
	result, err := database.Db.Exec("DELETE FROM task WHERE task_id = $1", id)
	if err != nil {
		http.Error(w, "Cann't remove task", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Task removed"})
}
