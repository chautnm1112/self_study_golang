package routes

import (
	"CS1_ToDoApp/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/tasks", controllers.GetAllTasks)
	r.GET("/tasks/:id", controllers.GetTaskByID)
	r.POST("/tasks", controllers.CreateNewTask)
	r.PUT("/tasks/:id", controllers.UpdateTask)
	r.DELETE("/tasks/:id", controllers.DeleteTask)

	return r
}
