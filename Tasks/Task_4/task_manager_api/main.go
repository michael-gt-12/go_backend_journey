package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"task_manager_api/data"
	"task_manager_api/router"
)

func main() {
	_ = godotenv.Load() // fine if .env is missing

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI is not set")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "taskmanager"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	client, err := data.Connect(uri, dbName)
	if err != nil {
		log.Fatalf("database unavailable: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = client.Disconnect(ctx)
	}()
	log.Println("connected to MongoDB")

	r := router.SetupRouter()
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
