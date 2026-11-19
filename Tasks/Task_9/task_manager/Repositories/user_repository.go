package Repositories

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"task_manager/Domain"
)

// UserRepository is the persistence contract for users.
// Implementations must return Domain.ErrNotFound when no user matches.
type UserRepository interface {
	Create(ctx context.Context, user Domain.User) (Domain.User, error)
	GetByUsername(ctx context.Context, username string) (Domain.User, error)
	GetByID(ctx context.Context, id string) (Domain.User, error)
	Count(ctx context.Context) (int64, error)
	UpdateRole(ctx context.Context, id string, role string) error
}

// ---- MongoDB implementation ----

// userDoc is the MongoDB representation of a user.
type userDoc struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"` // bcrypt hash
	Role     string             `bson:"role"`
}

func (d userDoc) toDomain() Domain.User {
	return Domain.User{
		ID:       d.ID.Hex(),
		Username: d.Username,
		Password: d.Password,
		Role:     d.Role,
	}
}

type mongoUserRepository struct {
	coll *mongo.Collection
}

// NewUserRepository returns a MongoDB-backed UserRepository.
func NewUserRepository(db *mongo.Database) UserRepository {
	return &mongoUserRepository{coll: db.Collection("users")}
}

// EnsureUserIndexes creates the unique index on username. Call once at startup.
func EnsureUserIndexes(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}

func (r *mongoUserRepository) Create(ctx context.Context, user Domain.User) (Domain.User, error) {
	d := userDoc{Username: user.Username, Password: user.Password, Role: user.Role}
	res, err := r.coll.InsertOne(ctx, d)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return Domain.User{}, Domain.ErrUsernameTaken
		}
		return Domain.User{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid.Hex()
	}
	return user, nil
}

func (r *mongoUserRepository) GetByUsername(ctx context.Context, username string) (Domain.User, error) {
	var d userDoc
	err := r.coll.FindOne(ctx, bson.M{"username": username}).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return Domain.User{}, Domain.ErrNotFound
	}
	if err != nil {
		return Domain.User{}, err
	}
	return d.toDomain(), nil
}

func (r *mongoUserRepository) GetByID(ctx context.Context, id string) (Domain.User, error) {
	oid, err := parseID(id)
	if err != nil {
		return Domain.User{}, err
	}
	var d userDoc
	err = r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return Domain.User{}, Domain.ErrNotFound
	}
	if err != nil {
		return Domain.User{}, err
	}
	return d.toDomain(), nil
}

func (r *mongoUserRepository) Count(ctx context.Context) (int64, error) {
	return r.coll.CountDocuments(ctx, bson.M{})
}

func (r *mongoUserRepository) UpdateRole(ctx context.Context, id string, role string) error {
	oid, err := parseID(id)
	if err != nil {
		return err
	}
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"role": role}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return Domain.ErrNotFound
	}
	return nil
}
