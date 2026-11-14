package Repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"task_manager/Domain"
)

// TaskRepository is the persistence contract for tasks.
// Implementations must return Domain.ErrNotFound when no task matches.
type TaskRepository interface {
	GetAll(ctx context.Context) ([]Domain.Task, error)
	GetByID(ctx context.Context, id string) (Domain.Task, error)
	Create(ctx context.Context, task Domain.Task) (Domain.Task, error)
	Update(ctx context.Context, id string, task Domain.Task) (Domain.Task, error)
	Delete(ctx context.Context, id string) error

	// ClaimDueTasks atomically moves up to limit tasks that are due on or before
	// dueBefore from pending to queued and returns only the tasks it moved.
	ClaimDueTasks(ctx context.Context, dueBefore time.Time, limit int) ([]Domain.Task, error)
	// SetReminderStatus sets reminder_status; returns Domain.ErrNotFound if no task matches.
	SetReminderStatus(ctx context.Context, id string, status string) error
}

// ---- MongoDB implementation ----

// taskDoc is the MongoDB representation of a task (keeps bson details out of Domain).
type taskDoc struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Title          string             `bson:"title"`
	Description    string             `bson:"description"`
	DueDate        time.Time          `bson:"due_date"`
	Status         string             `bson:"status"`
	ReminderStatus string             `bson:"reminder_status"`
}

func (d taskDoc) toDomain() Domain.Task {
	rs := d.ReminderStatus
	if rs == "" { // documents created before Task 8 have no field
		rs = Domain.ReminderPending
	}
	return Domain.Task{
		ID:             d.ID.Hex(),
		Title:          d.Title,
		Description:    d.Description,
		DueDate:        d.DueDate,
		Status:         d.Status,
		ReminderStatus: rs,
	}
}

// parseID converts a hex string to an ObjectID; a bad id is an invalid-input error.
func parseID(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("%w: invalid id", Domain.ErrInvalidInput)
	}
	return oid, nil
}

type mongoTaskRepository struct {
	coll *mongo.Collection
}

// NewTaskRepository returns a MongoDB-backed TaskRepository.
func NewTaskRepository(db *mongo.Database) TaskRepository {
	return &mongoTaskRepository{coll: db.Collection("tasks")}
}

func (r *mongoTaskRepository) GetAll(ctx context.Context) ([]Domain.Task, error) {
	cur, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var docs []taskDoc
	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}
	tasks := make([]Domain.Task, 0, len(docs))
	for _, d := range docs {
		tasks = append(tasks, d.toDomain())
	}
	return tasks, nil
}

func (r *mongoTaskRepository) GetByID(ctx context.Context, id string) (Domain.Task, error) {
	oid, err := parseID(id)
	if err != nil {
		return Domain.Task{}, err
	}
	var d taskDoc
	err = r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return Domain.Task{}, Domain.ErrNotFound
	}
	if err != nil {
		return Domain.Task{}, err
	}
	return d.toDomain(), nil
}

func (r *mongoTaskRepository) Create(ctx context.Context, task Domain.Task) (Domain.Task, error) {
	if task.ReminderStatus == "" {
		task.ReminderStatus = Domain.ReminderPending
	}
	d := taskDoc{
		Title:          task.Title,
		Description:    task.Description,
		DueDate:        task.DueDate,
		Status:         task.Status,
		ReminderStatus: task.ReminderStatus,
	}
	res, err := r.coll.InsertOne(ctx, d)
	if err != nil {
		return Domain.Task{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		task.ID = oid.Hex()
	}
	return task, nil
}

func (r *mongoTaskRepository) Update(ctx context.Context, id string, task Domain.Task) (Domain.Task, error) {
	oid, err := parseID(id)
	if err != nil {
		return Domain.Task{}, err
	}
	update := bson.M{"$set": bson.M{
		"title":           task.Title,
		"description":     task.Description,
		"due_date":        task.DueDate,
		"status":          task.Status,
		"reminder_status": task.ReminderStatus,
	}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var d taskDoc
	err = r.coll.FindOneAndUpdate(ctx, bson.M{"_id": oid}, update, opts).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return Domain.Task{}, Domain.ErrNotFound
	}
	if err != nil {
		return Domain.Task{}, err
	}
	return d.toDomain(), nil
}

func (r *mongoTaskRepository) Delete(ctx context.Context, id string) error {
	oid, err := parseID(id)
	if err != nil {
		return err
	}
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return Domain.ErrNotFound
	}
	return nil
}

func (r *mongoTaskRepository) ClaimDueTasks(ctx context.Context, dueBefore time.Time, limit int) ([]Domain.Task, error) {
	filter := bson.M{
		// $in with nil also matches documents that have no reminder_status yet.
		"reminder_status": bson.M{"$in": bson.A{Domain.ReminderPending, nil}},
		"due_date":        bson.M{"$lte": dueBefore}, // overdue tasks included
		"status":          bson.M{"$ne": "done"},
	}
	update := bson.M{"$set": bson.M{"reminder_status": Domain.ReminderQueued}}
	opts := options.FindOneAndUpdate().
		SetSort(bson.D{{Key: "due_date", Value: 1}}).
		SetReturnDocument(options.After)

	claimed := make([]Domain.Task, 0, max(limit, 0))
	for len(claimed) < limit {
		var d taskDoc
		err := r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&d)
		if errors.Is(err, mongo.ErrNoDocuments) {
			break // nothing left to claim
		}
		if err != nil {
			// Return what was already claimed so those tasks are not stranded in "queued".
			return claimed, err
		}
		claimed = append(claimed, d.toDomain())
	}
	return claimed, nil
}

func (r *mongoTaskRepository) SetReminderStatus(ctx context.Context, id string, status string) error {
	switch status {
	case Domain.ReminderPending, Domain.ReminderQueued, Domain.ReminderSent, Domain.ReminderFailed:
	default:
		return fmt.Errorf("%w: invalid reminder status %q", Domain.ErrInvalidInput, status)
	}
	oid, err := parseID(id)
	if err != nil {
		return err
	}
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"reminder_status": status}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return Domain.ErrNotFound
	}
	return nil
}
