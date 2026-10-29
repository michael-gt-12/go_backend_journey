package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task represents a task stored in MongoDB.
type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title" binding:"required"`
	Description string             `json:"description" bson:"description"`
	DueDate     time.Time          `json:"due_date" bson:"due_date"`
	Status      string             `json:"status" bson:"status" binding:"required,oneof=pending in_progress done"`
}
