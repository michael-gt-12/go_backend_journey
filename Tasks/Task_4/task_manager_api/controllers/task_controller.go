package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"task_manager_api/data"
	"task_manager_api/models"
)

// GetTasks handles GET /tasks.
func GetTasks(c *gin.Context) {
	tasks, err := data.GetAllTasks()
	if err != nil {
		log.Printf("GetAllTasks failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// GetTask handles GET /tasks/:id.
func GetTask(c *gin.Context) {
	task, err := data.GetTaskByID(c.Param("id"))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		case errors.Is(err, data.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		default:
			log.Printf("GetTaskByID failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch task"})
		}
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
		log.Printf("CreateTask failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// UpdateTask handles PUT /tasks/:id (implemented on Day 12).
func UpdateTask(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented yet"})
}

// DeleteTask handles DELETE /tasks/:id (implemented on Day 12).
func DeleteTask(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented yet"})
}
