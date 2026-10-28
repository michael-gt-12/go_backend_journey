package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Task struct {
	ID          primitive.ObjectID `bson:"id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	DueDate     time.Time          `bson:"due_date"`
	Status      string             `bson:"status"`
}

// give every database call a 5 second timeout
func newCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func main() {
	_ = godotenv.Load()
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI is not set")
	}

	// create a connection with the database
	ctx, cancel := newCtx()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	cancel()
	if err != nil {
		log.Fatal("connect: ", err)
	}

	defer func() {
		ctx, cancel = newCtx()
		defer cancel()
		_ = client.Disconnect(ctx)
	}()

	// PING the server to check if connection is established
	ctx, cancel = newCtx()
	err = client.Ping(ctx, nil)
	cancel()
	if err != nil {
		log.Fatal("ping:", err)
	}
	fmt.Println("connected to MongoDB")

	coll := client.Database("taskmanager").Collection("tasks")

	id := insertTask(coll)
	findOne(coll, id)
	findMany(coll)
	updateStatus(coll, id)
	deleteTask(coll, id)
}

// Insert in to the DB
func insertTask(coll *mongo.Collection) primitive.ObjectID {
	ctx, cancel := newCtx()
	defer cancel()
	task := Task{
		Title:       "Write Day 10 notes",
		Description: "Practice the Go driver",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      "pending",
	}
	res, err := coll.InsertOne(ctx, task)
	if err != nil {
		log.Fatal("insert: ", err)
	}
	id := res.InsertedID.(primitive.ObjectID)
	fmt.Println("inserted with id:", id.Hex())
	return id
}

// 4) FIND ONE by _id
func findOne(coll *mongo.Collection, id primitive.ObjectID) {
	ctx, cancel := newCtx()
	defer cancel()

	var found Task
	err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&found)
	if errors.Is(err, mongo.ErrNoDocuments) {
		fmt.Println("no task with that id")
		return
	}
	if err != nil {
		log.Fatal("findone:", err)
	}
	fmt.Printf("found: %+v\n", found)
}

// 5) FIND with a filter (many documents)
func findMany(coll *mongo.Collection) {
	ctx, cancel := newCtx()
	defer cancel()

	cursor, err := coll.Find(ctx, bson.M{"status": "pending"})
	if err != nil {
		log.Fatal("find:", err)
	}
	defer cursor.Close(ctx) // free the cursor

	var tasks []Task
	if err := cursor.All(ctx, &tasks); err != nil {
		log.Fatal("cursor:", err)
	}
	fmt.Println("pending tasks:", len(tasks))
	for _, t := range tasks {
		fmt.Println(" -", t.Title)
	}
}

// 6) UPDATE ONE with $set
func updateStatus(coll *mongo.Collection, id primitive.ObjectID) {
	ctx, cancel := newCtx()
	defer cancel()

	res, err := coll.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"status": "done"}},
	)
	if err != nil {
		log.Fatal("update:", err)
	}
	fmt.Printf("matched=%d modified=%d\n", res.MatchedCount, res.ModifiedCount)
}

// 7) DELETE ONE
func deleteTask(coll *mongo.Collection, id primitive.ObjectID) {
	ctx, cancel := newCtx()
	defer cancel()

	res, err := coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		log.Fatal("delete:", err)
	}
	fmt.Println("deleted:", res.DeletedCount)
}
