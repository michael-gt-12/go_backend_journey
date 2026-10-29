package router

import (
	"github.com/gin-gonic/gin"

	"task_manager_api/controllers"
)

// SetupRouter registers all task routes.
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/tasks", controllers.GetTasks)
	r.GET("/tasks/:id", controllers.GetTask)
	r.POST("/tasks", controllers.CreateTask)
	r.PUT("/tasks/:id", controllers.UpdateTask)
	r.DELETE("/tasks/:id", controllers.DeleteTask)

	return r
}
