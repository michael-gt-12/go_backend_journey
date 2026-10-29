package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"task_manager_api/models"
)

// Sentinel errors so controllers can map them to HTTP status codes.
var (
	ErrInvalidID = errors.New("invalid task id")
	ErrNotFound  = errors.New("task not found")
)

const opTimeout = 5 * time.Second

// Collection is the MongoDB collection used for tasks.
// It is set once by Connect at startup.
var Collection *mongo.Collection

// Connect opens a MongoDB connection, verifies it with Ping, and
// stores the tasks collection. It returns an error if the database
// is unreachable so main can fail fast.
func Connect(uri, dbName string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping: %w", err)
	}

	Collection = client.Database(dbName).Collection("tasks")
	return client, nil
}

// GetAllTasks returns every task. It always returns a non-nil slice
// so the API responds with [] instead of null.
func GetAllTasks() ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()

	cursor, err := Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	if tasks == nil {
		tasks = []models.Task{}
	}
	return tasks, nil
}

// GetTaskByID returns one task.
// It returns ErrInvalidID for a malformed id and ErrNotFound when no document matches.
func GetTaskByID(id string) (models.Task, error) {
	var task models.Task

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return task, ErrInvalidID
	}

	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()

	err = Collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&task)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return task, ErrNotFound
		}
		return task, err
	}
	return task, nil
}

// CreateTask inserts a new task and returns it with its generated id.
// Any id sent by the client is ignored.
func CreateTask(task models.Task) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()

	task.ID = primitive.NilObjectID // let MongoDB generate the id

	result, err := Collection.InsertOne(ctx, task)
	if err != nil {
		return models.Task{}, err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return models.Task{}, errors.New("unexpected inserted id type")
	}
	task.ID = oid
	return task, nil
}
