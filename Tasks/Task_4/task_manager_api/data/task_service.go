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
	ErrInvalidID     = errors.New("invalid task id")
	ErrNotFound      = errors.New("task not found")
	ErrEmptyUpdate   = errors.New("no fields to update")
	ErrTimeout       = errors.New("database operation timed out")
	ErrDBUnavailable = errors.New("database unavailable")
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

// wrapDBError classifies a driver error so the controller can tell
// a timeout or an unreachable database apart from other failures.
// Other errors are returned unchanged (the controller logs them and
// answers with a generic 500).
func wrapDBError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, context.DeadlineExceeded), mongo.IsTimeout(err):
		return fmt.Errorf("%w: %v", ErrTimeout, err)
	case mongo.IsNetworkError(err):
		return fmt.Errorf("%w: %v", ErrDBUnavailable, err)
	default:
		return err
	}
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

// UpdateTask replaces the editable fields of the task with the given id
// and returns the updated task.
// It returns ErrInvalidID for a malformed id and ErrNotFound when
// MatchedCount is 0 (no document has that id). Any id in the body is ignored.
func UpdateTask(id string, task models.Task) (models.Task, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Task{}, ErrInvalidID
	}

	// Build the $set document from the struct using its own bson tags,
	// then drop _id: MongoDB does not allow changing it.
	raw, err := bson.Marshal(task)
	if err != nil {
		return models.Task{}, err
	}
	var fields bson.M
	if err := bson.Unmarshal(raw, &fields); err != nil {
		return models.Task{}, err
	}
	delete(fields, "_id")
	if len(fields) == 0 {
		return models.Task{}, ErrEmptyUpdate
	}

	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()

	result, err := Collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": fields})
	if err != nil {
		return models.Task{}, wrapDBError(err)
	}
	// Use MatchedCount, not ModifiedCount: sending identical values
	// modifies nothing but is still a successful update.
	if result.MatchedCount == 0 {
		return models.Task{}, ErrNotFound
	}

	return GetTaskByID(id)
}

// DeleteTask removes the task with the given id.
// It returns ErrInvalidID for a malformed id and ErrNotFound when
// DeletedCount is 0 (no document has that id).
func DeleteTask(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}

	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()

	result, err := Collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return wrapDBError(err)
	}
	if result.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}
