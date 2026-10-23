package data

import (
	"errors"
	"sort"
	"strconv"
	"sync"
	"task_manager_api/models"
	"time"
)

// ErrTaskNotFound is returned when no task matches the given id.
var ErrTaskNotFound = errors.New("task not found")

// mu protects tasks and nextID from concurrent requests.
var mu sync.RWMutex

// nextID is the id given to the next created task.
var nextID = 4

// tasks stores the tasks in memory, keyed by id.
var tasks = map[string]models.Task{
	"1": {ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "pending"},
	"2": {ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "in_progress"},
	"3": {ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "done"},
}

// GetAllTasks returns all tasks ordered by due date.
func GetAllTasks() []models.Task {
	mu.RLock()
	defer mu.RUnlock()

	result := []models.Task{}
	for _, task := range tasks {
		result = append(result, task)
	}
	// map order is random in Go, so sort to get a stable result
	sort.Slice(result, func(i, j int) bool {
		return result[i].DueDate.Before(result[j].DueDate)
	})
	return result
}

// GetTaskByID returns the task with the given id, or ErrTaskNotFound.
func GetTaskByID(id string) (models.Task, error) {
	mu.RLock()
	defer mu.RUnlock()

	task, found := tasks[id]
	if !found {
		return models.Task{}, ErrTaskNotFound
	}
	return task, nil
}

// CreateTask gives the task a new id, stores it and returns it.
func CreateTask(task models.Task) models.Task {
	mu.Lock()
	defer mu.Unlock()

	task.ID = strconv.Itoa(nextID)
	nextID++
	tasks[task.ID] = task
	return task
}

// UpdateTask changes only the non-empty fields of the task with the given id.
func UpdateTask(id string, updatedTask models.Task) (models.Task, error) {
	mu.Lock()
	defer mu.Unlock()

	task, found := tasks[id]
	if !found {
		return models.Task{}, ErrTaskNotFound
	}
	if updatedTask.Title != "" {
		task.Title = updatedTask.Title
	}
	if updatedTask.Description != "" {
		task.Description = updatedTask.Description
	}
	if updatedTask.Status != "" {
		task.Status = updatedTask.Status
	}
	if !updatedTask.DueDate.IsZero() {
		task.DueDate = updatedTask.DueDate
	}
	// a map stores copies, so save the changed task back
	tasks[id] = task
	return task, nil
}

// DeleteTaskByID removes the task with the given id, or returns ErrTaskNotFound.
func DeleteTaskByID(id string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, found := tasks[id]; !found {
		return ErrTaskNotFound
	}
	delete(tasks, id)
	return nil
}
