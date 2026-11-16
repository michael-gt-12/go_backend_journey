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
	"task_manager/Delivery/workers"
	"task_manager/Domain"
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

	reminderCfg, err := LoadReminderConfig()
	if err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

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

	notifier := Infrastructure.NewLogNotifier()
	reminderUC := Usecases.NewReminderUsecase(taskRepo, notifier, 3, 2*time.Second)

	controller := controllers.NewController(taskUC, userUC)
	router := routers.SetupRouter(controller, jwtService)

	jobs := make(chan Domain.Task, reminderCfg.QueueSize)
	scheduler := workers.NewScheduler(reminderUC, reminderCfg.Interval, reminderCfg.Window, jobs)
	pool := workers.NewWorkerPool(reminderCfg.Workers, reminderUC, 5*time.Second)

	// Temporary for today: a plain cancelable context. not  using signal.NotifyContext yet,
	// because router.Run blocks and would swallow Ctrl+C. Day 27 replaces this with
	// real graceful shutdown. appCtx is separate from the 10s startup ctx above.
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	pool.Start(jobs)
	go scheduler.Run(appCtx)

	log.Printf("server listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		appCancel()
		pool.Wait()
		_ = client.Disconnect(context.Background())
		log.Fatalf("server: %v", err)
	}
}
