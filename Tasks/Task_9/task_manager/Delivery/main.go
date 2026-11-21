package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}
	port := getenv("PORT", "8080")

	// Which database stores users AND tasks. Default is postgres (per the Task 9 spec).
	dbDriver := getenv("DB_DRIVER", "postgres")
	if dbDriver != "postgres" && dbDriver != "mongo" {
		log.Fatalf("invalid DB_DRIVER %q (use postgres or mongo)", dbDriver)
	}

	reminderCfg, err := LoadReminderConfig()
	if err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	// Pick the repositories. This switch is the only place that knows which database is used.
	var (
		taskRepo Repositories.TaskRepository
		userRepo Repositories.UserRepository
		pgPool   *pgxpool.Pool
		client   *mongo.Client
	)

	switch dbDriver {
	case "postgres":
		// The helper has its own 5s ping timeout, so a plain background context is fine.
		pgPool, err = Infrastructure.NewPostgresPool(context.Background())
		if err != nil {
			log.Fatalf("postgres: %v", err)
		}
		userRepo = Repositories.NewPostgresUserRepository(pgPool)
		taskRepo = Repositories.NewPostgresTaskRepository(pgPool)
		log.Println("storage: using PostgreSQL for users and tasks")

	case "mongo":
		mongoURI := os.Getenv("MONGO_URI")
		if mongoURI == "" {
			log.Fatal("MONGO_URI must be set when DB_DRIVER=mongo")
		}
		dbName := getenv("MONGO_DB", "taskmanager")

		// Connect and fail fast if MongoDB is unreachable.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
		if err != nil {
			cancel()
			log.Fatalf("mongo connect: %v", err)
		}
		if err := client.Ping(ctx, nil); err != nil {
			cancel()
			log.Fatalf("mongo ping: %v", err)
		}
		db := client.Database(dbName)

		// The unique index on username is only needed for Mongo;
		// in Postgres the UNIQUE constraint from the migration does this job.
		if err := Repositories.EnsureUserIndexes(ctx, db); err != nil {
			cancel()
			log.Fatalf("ensure indexes: %v", err)
		}
		cancel()

		userRepo = Repositories.NewMongoUserRepository(db)
		taskRepo = Repositories.NewTaskRepository(db)
		log.Println("storage: using MongoDB for users and tasks")
	}

	// Manual dependency injection: infrastructure -> repositories -> usecases -> controller -> router.
	passwordService := Infrastructure.NewPasswordService()
	jwtService := Infrastructure.NewJWTService(jwtSecret, 24*time.Hour)

	taskUC := Usecases.NewTaskUsecase(taskRepo)
	userUC := Usecases.NewUserUsecase(userRepo, passwordService, jwtService)

	notifier := Infrastructure.NewLogNotifier()
	reminderUC := Usecases.NewReminderUsecase(taskRepo, notifier, 3, 2*time.Second)

	controller := controllers.NewController(taskUC, userUC)
	router := routers.SetupRouter(controller, jwtService)

	jobs := make(chan Domain.Task, reminderCfg.QueueSize)
	scheduler := workers.NewScheduler(reminderUC, reminderCfg.Interval, reminderCfg.Window, jobs)
	pool := workers.NewWorkerPool(reminderCfg.Workers, reminderUC, 5*time.Second)

	// Ctrl+C / SIGTERM arrives here.
	sigCtx, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignals()

	// The scheduler gets its OWN context (not sigCtx) so we control the order:
	// HTTP first, then scheduler, then workers.
	schedCtx, cancelSched := context.WithCancel(context.Background())
	defer cancelSched()

	pool.Start(jobs)
	schedDone := make(chan struct{})
	go func() {
		defer close(schedDone)
		scheduler.Run(schedCtx)
	}()

	// Manual http.Server instead of router.Run, so we can call Shutdown.
	srv := &http.Server{Addr: ":" + port, Handler: router}
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("server listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// Wait for a signal or a server failure.
	select {
	case <-sigCtx.Done():
		log.Println("shutdown: signal received")
	case err := <-serverErr:
		log.Printf("server error: %v", err)
	}

	// 1. Stop accepting HTTP requests (10s timeout).
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: http error: %v", err)
	}
	log.Println("shutdown: http server stopped")

	// 2. Stop the scheduler; it closes the jobs channel on return.
	cancelSched()
	<-schedDone
	log.Println("shutdown: scheduler stopped, jobs channel closed")

	// 3. Wait for workers to finish in-flight jobs.
	pool.Wait()
	log.Println("shutdown: workers finished")

	// 4. Close the database last, after nothing can use it anymore.
	if pgPool != nil {
		pgPool.Close()
		log.Println("shutdown: postgres pool closed")
	}
	if client != nil {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("shutdown: mongo disconnect: %v", err)
		}
		log.Println("shutdown: mongo disconnected")
	}
	log.Println("shutdown: bye")
}
