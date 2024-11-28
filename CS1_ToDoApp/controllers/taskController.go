package controllers

import (
	"CS1_ToDoApp/database"
	"CS1_ToDoApp/models"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "log"
	"net/http"
	"time"
	_ "time"
)

func GetAllTasks(c *gin.Context) {
	var tasks []models.Task
	rows, err := database.Db.Query("SELECT task_id, task, completed, created_at, updated_at FROM task")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cann't get all task"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Task, &task.Completed, &task.CreatedAt, &task.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when getting all task"})
			return
		}
		tasks = append(tasks, task)
	}

	c.JSON(http.StatusOK, tasks)

}

func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	row := database.Db.QueryRow("SELECT task_id, task, completed, created_at, updated_at FROM task WHERE task_id = $1", id)
	if err := row.Scan(&task.ID, &task.Task, &task.Completed, &task.CreatedAt, &task.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying task", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

func CreateNewTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	var lastInsertID int64
	err := database.Db.QueryRow(`INSERT INTO task(task, completed, created_at, updated_at)
		VALUES($1, $2, $3, $4) RETURNING task_id`, task.Task, task.Completed, task.CreatedAt, task.UpdatedAt).Scan(&lastInsertID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cann't create task"})
		return
	}

	task.ID = lastInsertID
	c.JSON(http.StatusCreated, task)

}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := database.Db.QueryRow("SELECT task_id, task, completed, created_at, updated_at FROM task WHERE task_id = $1", id).Scan(&task.ID, &task.Task, &task.Completed, &task.CreatedAt, &task.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying task", "details": err.Error()})
		}
		return
	}

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	task.UpdatedAt = time.Now()

	_, err := database.Db.Exec(`UPDATE task SET task = $1, completed = $2, updated_at = $3 WHERE task_id = $4`,
		task.Task, task.Completed, task.UpdatedAt, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cann't update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	result, err := database.Db.Exec("DELETE FROM task WHERE task_id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot remove task"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task removed"})
}
