package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"task_manager/Domain"
	"task_manager/Repositories"
	"task_manager/Usecases"
)

// Keys under which the auth middleware stores the authenticated user in the Gin context.
// IMPORTANT: they must match the c.Set(...) calls in Infrastructure/auth_middleWare.go.
const (
	ctxUserIDKey = "userID"
	ctxRoleKey   = "role"
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

// callerContext returns a request context that carries the authenticated caller.
// It writes a 401 and returns ok=false when the auth middleware did not set a user.
func callerContext(c *gin.Context) (context.Context, bool) {
	id := c.GetString(ctxUserIDKey)
	role := c.GetString(ctxRoleKey)
	if id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return nil, false
	}
	return Usecases.WithCaller(c.Request.Context(), Usecases.Caller{ID: id, Role: role}), true
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

// parseListQuery validates ?page=, ?limit= and ?status= (defaults: page 1, limit 10).
func parseListQuery(c *gin.Context) (Repositories.TaskFilter, error) {
	f := Repositories.TaskFilter{Page: 1, Limit: 10}

	if v := c.Query("page"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			return f, fmt.Errorf("%w: page must be an integer >= 1", Domain.ErrInvalidInput)
		}
		f.Page = n
	}
	if v := c.Query("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 100 {
			return f, fmt.Errorf("%w: limit must be an integer between 1 and 100", Domain.ErrInvalidInput)
		}
		f.Limit = n
	}
	if v := c.Query("status"); v != "" {
		if v != "pending" && v != "in_progress" && v != "done" {
			return f, fmt.Errorf("%w: status must be pending, in_progress or done", Domain.ErrInvalidInput)
		}
		f.Status = v
	}
	return f, nil
}

// GetTasks handles GET /tasks. The body stays an array (backward compatible);
// totals are returned in the X-Total-Count, X-Page and X-Limit headers.
func (h *Controller) GetTasks(c *gin.Context) {
	ctx, ok := callerContext(c)
	if !ok {
		return
	}
	filter, err := parseListQuery(c)
	if err != nil {
		handleError(c, err)
		return
	}
	tasks, total, err := h.tasks.ListTasks(ctx, filter)
	if err != nil {
		handleError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.Itoa(total))
	c.Header("X-Page", strconv.Itoa(filter.Page))
	c.Header("X-Limit", strconv.Itoa(filter.Limit))
	c.JSON(http.StatusOK, tasks)
}

// GetTask handles GET /tasks/:id.
func (h *Controller) GetTask(c *gin.Context) {
	ctx, ok := callerContext(c)
	if !ok {
		return
	}
	task, err := h.tasks.GetTaskByID(ctx, c.Param("id"))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

// CreateTask handles POST /tasks. The owner is taken from the token, never from the body.
func (h *Controller) CreateTask(c *gin.Context) {
	ctx, ok := callerContext(c)
	if !ok {
		return
	}
	var in Domain.Task
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}
	in.ID = ""
	in.OwnerID = ""
	task, err := h.tasks.CreateTask(ctx, in)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, task)
}

// UpdateTask handles PUT /tasks/:id (owner or admin).
func (h *Controller) UpdateTask(c *gin.Context) {
	ctx, ok := callerContext(c)
	if !ok {
		return
	}
	var in Domain.Task
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}
	in.OwnerID = ""
	task, err := h.tasks.UpdateTask(ctx, c.Param("id"), in)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

// DeleteTask handles DELETE /tasks/:id (owner or admin).
func (h *Controller) DeleteTask(c *gin.Context) {
	ctx, ok := callerContext(c)
	if !ok {
		return
	}
	if err := h.tasks.DeleteTask(ctx, c.Param("id")); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}
