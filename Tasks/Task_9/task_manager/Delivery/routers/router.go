package routers

import (
	"github.com/gin-gonic/gin"

	"task_manager/Delivery/controllers"
	"task_manager/Infrastructure"
)

// SetupRouter registers all routes and attaches the auth middleware.
func SetupRouter(h *controllers.Controller, jwtService *Infrastructure.JWTService) *gin.Engine {
	r := gin.Default()

	// Public routes
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)

	// Any logged-in user
	auth := r.Group("/")
	auth.Use(Infrastructure.AuthMiddleware(jwtService))
	auth.GET("/tasks", h.GetTasks)
	auth.GET("/tasks/:id", h.GetTask)

	// Admin only
	admin := auth.Group("/")
	admin.Use(Infrastructure.AdminOnly())
	admin.POST("/tasks", h.CreateTask)
	admin.PUT("/tasks/:id", h.UpdateTask)
	admin.DELETE("/tasks/:id", h.DeleteTask)
	admin.POST("/promote/:id", h.Promote)

	return r
}
