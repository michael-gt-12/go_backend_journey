package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"task_manager/Delivery/controllers"
	"task_manager/Delivery/routers"
	"task_manager/Infrastructure"
	"task_manager/Repositories"
	"task_manager/Usecases"
)

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	// .env is optional; real environment variables win.
	_ = godotenv.Load()

	mongoURI := os.Getenv("MONGO_URI")
	jwtSecret := os.Getenv("JWT_SECRET")
	if mongoURI == "" || jwtSecret == "" {
		log.Fatal("MONGO_URI and JWT_SECRET must be set")
	}
	dbName := getenv("MONGO_DB", "taskmanager")
	port := getenv("PORT", "8080")

	// Connect to MongoDB and fail fast if it is unreachable.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("mongo ping: %v", err)
	}
	db := client.Database(dbName)

	if err := Repositories.EnsureUserIndexes(ctx, db); err != nil {
		log.Fatalf("ensure indexes: %v", err)
	}

	// Manual dependency injection: infrastructure -> repositories -> usecases -> controller -> router.
	passwordService := Infrastructure.NewPasswordService()
	jwtService := Infrastructure.NewJWTService(jwtSecret, 24*time.Hour)

	taskRepo := Repositories.NewTaskRepository(db)
	userRepo := Repositories.NewUserRepository(db)

	taskUC := Usecases.NewTaskUsecase(taskRepo)
	userUC := Usecases.NewUserUsecase(userRepo, passwordService, jwtService)

	controller := controllers.NewController(taskUC, userUC)
	router := routers.SetupRouter(controller, jwtService)

	log.Printf("server listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		_ = client.Disconnect(context.Background())
		log.Fatalf("server: %v", err)
	}
}
