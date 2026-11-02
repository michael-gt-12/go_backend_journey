package data

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"task_manager_api/models"
)

// Sentinel errors so controllers can map them to HTTP status codes.
var (
	ErrUsernameTaken      = errors.New("username already taken")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrInvalidUserID      = errors.New("invalid user id")
	ErrUserNotFound       = errors.New("user not found")
)

// UsersCollection is the MongoDB collection for users. Set by Connect.
var UsersCollection *mongo.Collection

// dummyHash is compared against when the username does not exist, so a
// "no such user" login takes about as long as a "wrong password" login.
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("dummy-password"), bcrypt.DefaultCost)

// EnsureUserIndexes creates the unique index on username.
// This index is what really guarantees uniqueness, even under concurrent requests.
func EnsureUserIndexes(ctx context.Context) error {
	_, err := UsersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}

// Register hashes the password and stores a new user.
// The first registered user becomes admin, everyone after is a normal user.
// It returns ErrUsernameTaken when the unique index rejects the insert.
func Register(username, password string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	role := models.RoleUser
	count, err := UsersCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return models.User{}, wrapDBError(err)
	}
	if count == 0 {
		role = models.RoleAdmin
	}

	user := models.User{Username: username, Password: string(hash), Role: role}
	result, err := UsersCollection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return models.User{}, ErrUsernameTaken
		}
		return models.User{}, wrapDBError(err)
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

// Authenticate checks the credentials and returns the user.
// Unknown username and wrong password return the SAME error, so a
// caller cannot discover which usernames exist.
func Authenticate(username, password string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()

	var user models.User
	err := UsersCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(password))
			return models.User{}, ErrInvalidCredentials
		}
		return models.User{}, wrapDBError(err)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return models.User{}, ErrInvalidCredentials
	}
	return user, nil
}

// PromoteUser sets the role of the user with the given id to admin.
// It returns ErrInvalidUserID for a malformed id and ErrUserNotFound
// when MatchedCount is 0.
func PromoteUser(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidUserID
	}

	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()

	result, err := UsersCollection.UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{"role": models.RoleAdmin}})
	if err != nil {
		return wrapDBError(err)
	}
	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}
