package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"task_manager/Domain"
	"task_manager/Usecases"
)

// Controller holds the usecases the HTTP handlers depend on (interfaces only).
type Controller struct {
	tasks Usecases.TaskUsecase
	users Usecases.UserUsecase
}

// NewController builds a Controller with injected usecases.
func NewController(t Usecases.TaskUsecase, u Usecases.UserUsecase) *Controller {
	return &Controller{tasks: t, users: u}
}

type credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// handleError maps domain errors to HTTP status codes.
func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, Domain.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, Domain.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case errors.Is(err, Domain.ErrUsernameTaken):
		c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
	case errors.Is(err, Domain.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
	default:
		log.Printf("internal error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

// ---- auth ----

// Register handles POST /register.
func (h *Controller) Register(c *gin.Context) {
	var in credentials
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}
	user, err := h.users.Register(c.Request.Context(), in.Username, in.Password)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, user)
}

// Login handles POST /login.
func (h *Controller) Login(c *gin.Context) {
	var in credentials
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}
	token, err := h.users.Login(c.Request.Context(), in.Username, in.Password)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Promote handles POST /promote/:id (admin only).
func (h *Controller) Promote(c *gin.Context) {
	if err := h.users.Promote(c.Request.Context(), c.Param("id")); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user promoted to admin"})
}

// ---- tasks ----

// GetTasks handles GET /tasks.
func (h *Controller) GetTasks(c *gin.Context) {
	tasks, err := h.tasks.GetAllTasks(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// GetTask handles GET /tasks/:id.
func (h *Controller) GetTask(c *gin.Context) {
	task, err := h.tasks.GetTaskByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

// CreateTask handles POST /tasks (admin only).
func (h *Controller) CreateTask(c *gin.Context) {
	var in Domain.Task
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}
	in.ID = ""
	task, err := h.tasks.CreateTask(c.Request.Context(), in)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, task)
}

// UpdateTask handles PUT /tasks/:id (admin only).
func (h *Controller) UpdateTask(c *gin.Context) {
	var in Domain.Task
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}
	task, err := h.tasks.UpdateTask(c.Request.Context(), c.Param("id"), in)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

// DeleteTask handles DELETE /tasks/:id (admin only).
func (h *Controller) DeleteTask(c *gin.Context) {
	if err := h.tasks.DeleteTask(c.Request.Context(), c.Param("id")); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}
