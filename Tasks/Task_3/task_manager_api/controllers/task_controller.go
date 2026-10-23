package controllers

import (
	"errors"
	"net/http"
	"task_manager_api/data"
	"task_manager_api/models"

	"github.com/gin-gonic/gin"
)

// GetTasks handles GET /tasks.
func GetTasks(c *gin.Context) {
	c.JSON(http.StatusOK, data.GetAllTasks())
}

// GetTask handles GET /tasks/:id.
func GetTask(c *gin.Context) {
	task, err := data.GetTaskByID(c.Param("id"))
	if errors.Is(err, data.ErrTaskNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// CreateTask handles POST /tasks.
func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, data.CreateTask(task))
}

// UpdateTask handles PUT /tasks/:id.
func UpdateTask(c *gin.Context) {
	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task, err := data.UpdateTask(c.Param("id"), updatedTask)
	if errors.Is(err, data.ErrTaskNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// DeleteTask handles DELETE /tasks/:id.
func DeleteTask(c *gin.Context) {
	err := data.DeleteTaskByID(c.Param("id"))
	if errors.Is(err, data.ErrTaskNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task removed"})
}
