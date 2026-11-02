package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"task_manager_api/data"
	"task_manager_api/middleware"
	"task_manager_api/models"
)

// respondError maps a data-layer error to an HTTP status code.
// Client mistakes are 400/401/404/409, timeouts and an unreachable
// database are 504/503, and anything else is logged and answered with
// a generic 500 so no driver details leak to the client.
func respondError(c *gin.Context, err error, fallbackMsg string) {
	switch {
	case errors.Is(err, data.ErrInvalidID):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
	case errors.Is(err, data.ErrInvalidUserID):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
	case errors.Is(err, data.ErrEmptyUpdate):
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
	case errors.Is(err, data.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
	case errors.Is(err, data.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	case errors.Is(err, data.ErrUsernameTaken):
		c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
	case errors.Is(err, data.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
	case errors.Is(err, data.ErrTimeout):
		log.Printf("%s: %v", fallbackMsg, err)
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "database timed out, please try again"})
	case errors.Is(err, data.ErrDBUnavailable):
		log.Printf("%s: %v", fallbackMsg, err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "database unavailable, please try again later"})
	default:
		log.Printf("%s: %v", fallbackMsg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fallbackMsg})
	}
}

// ---------- Auth ----------

// Register handles POST /register.
func Register(c *gin.Context) {
	var input models.Credentials
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username (3-32 characters) and password (8-72 characters) are required",
		})
		return
	}

	user, err := data.Register(input.Username, input.Password)
	if err != nil {
		respondError(c, err, "failed to register user")
		return
	}
	// models.User hides the password hash with json:"-".
	c.JSON(http.StatusCreated, user)
}

// Login handles POST /login.
// Unknown username and wrong password both produce the same 401 message.
func Login(c *gin.Context) {
	var input models.Credentials
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username and password are required",
		})
		return
	}

	user, err := data.Authenticate(input.Username, input.Password)
	if err != nil {
		respondError(c, err, "failed to log in")
		return
	}

	token, err := middleware.GenerateToken(user)
	if err != nil {
		log.Printf("GenerateToken failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to log in"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// PromoteUser handles POST /promote/:id (admin only).
func PromoteUser(c *gin.Context) {
	if err := data.PromoteUser(c.Param("id")); err != nil {
		respondError(c, err, "failed to promote user")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user promoted to admin"})
}

// ---------- Tasks ----------

// GetTasks handles GET /tasks.
func GetTasks(c *gin.Context) {
	tasks, err := data.GetAllTasks()
	if err != nil {
		respondError(c, err, "failed to fetch tasks")
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// GetTask handles GET /tasks/:id.
func GetTask(c *gin.Context) {
	task, err := data.GetTaskByID(c.Param("id"))
	if err != nil {
		respondError(c, err, "failed to fetch task")
		return
	}
	c.JSON(http.StatusOK, task)
}

// CreateTask handles POST /tasks.
func CreateTask(c *gin.Context) {
	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := data.CreateTask(input)
	if err != nil {
		respondError(c, err, "failed to create task")
		return
	}
	c.JSON(http.StatusCreated, created)
}

// UpdateTask handles PUT /tasks/:id.
func UpdateTask(c *gin.Context) {
	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := data.UpdateTask(c.Param("id"), input)
	if err != nil {
		respondError(c, err, "failed to update task")
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteTask handles DELETE /tasks/:id.
func DeleteTask(c *gin.Context) {
	if err := data.DeleteTask(c.Param("id")); err != nil {
		respondError(c, err, "failed to delete task")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}
