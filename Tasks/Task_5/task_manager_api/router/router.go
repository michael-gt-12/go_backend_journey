package router

import (
	"github.com/gin-gonic/gin"

	"task_manager_api/controllers"
	"task_manager_api/middleware"
)

// SetupRouter registers all routes.
//
//	public:  POST /register, POST /login
//	any logged-in user: GET /tasks, GET /tasks/:id
//	admin only: POST/PUT/DELETE on tasks, POST /promote/:id
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Public routes.
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Routes for any authenticated user.
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	auth.GET("/tasks", controllers.GetTasks)
	auth.GET("/tasks/:id", controllers.GetTask)

	// Routes for admins only (AuthRequired runs first, then AdminOnly).
	admin := auth.Group("/")
	admin.Use(middleware.AdminOnly())
	admin.POST("/tasks", controllers.CreateTask)
	admin.PUT("/tasks/:id", controllers.UpdateTask)
	admin.DELETE("/tasks/:id", controllers.DeleteTask)
	admin.POST("/promote/:id", controllers.PromoteUser)

	return r
}
