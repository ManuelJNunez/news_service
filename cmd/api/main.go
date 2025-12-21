package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ManuelJNunez/news_service/internal/config"
	"github.com/ManuelJNunez/news_service/internal/health"
	"github.com/ManuelJNunez/news_service/internal/news"
	"github.com/ManuelJNunez/news_service/internal/user"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	// 1) Load configuration
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2) Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	logger.Info("starting service", slog.String("port", cfg.HTTPPort))

	// 3) Connect to PostgreSQL database
	db, err := initDB(cfg, logger)

	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("failed to close database", slog.Any("error", err))
		}
	}()

	// 4) Connect to MongoDB
	mongoClient, err := initMongoDB(cfg, logger)
	if err != nil {
		logger.Error("failed to connect to mongodb", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			logger.Error("failed to disconnect from mongodb", slog.Any("error", err))
		}
	}()

	// 5) Build dependencies from news domain
	newsRepo := news.NewRepository(db)
	newsSvc := news.NewService(newsRepo)
	newsHandler := news.NewHandler(newsSvc)

	// 6) Build dependencies from user domain
	usersCollection := mongoClient.Database("app").Collection("users")
	userRepo := user.NewRepository(usersCollection)
	userSvc := user.NewService(userRepo)
	userHandler := user.NewHandler(userSvc)

	// 7) Configure Gin (web framework)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router_group := router.Group("")

	// Register health route
	healthHandler := health.NewHandler()
	health.RegisterRoutes(router_group, healthHandler)

	// Register news routes
	news.RegisterRoutes(router_group, newsHandler)

	// Register user routes
	user.RegisterRoutes(router_group, userHandler)

	// 8) Configure HTTP server
	addr := ":" + cfg.HTTPPort
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// 9) Launch webserver on a separate Thread using a Goroutine
	go func() {
		logger.Info("HTTP server listening", slog.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server listen error", slog.Any("error", err))
		}
	}()

	// 10) Configure signal handling and pause execution until SIGINT or SIGTERM is received
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", slog.Any("error", err))
	}

	logger.Info("server exiting")
}

func initDB(cfg *config.Config, logger *slog.Logger) (*sql.DB, error) {
	// Open DB driver connection (using the postgres driver)
	db, err := sql.Open("postgres", cfg.DB_DSN)

	if err != nil {
		return nil, err
	}

	// Initialize context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check and open DB connection
	if err := db.PingContext(ctx); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			logger.Error("failed to close database after ping error", slog.Any("error", closeErr))
		}
		logger.Error("database ping failed", slog.Any("error", err))
		return nil, err
	}

	logger.Info("database connection established")
	return db, nil
}

func initMongoDB(cfg *config.Config, logger *slog.Logger) (*mongo.Client, error) {
	// Initialize context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoDB_URI))
	if err != nil {
		return nil, err
	}

	// Ping MongoDB to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		if disconnectErr := client.Disconnect(ctx); disconnectErr != nil {
			logger.Error("failed to disconnect from mongodb after ping error", slog.Any("error", disconnectErr))
		}
		logger.Error("mongodb ping failed", slog.Any("error", err))
		return nil, err
	}

	logger.Info("mongodb connection successfully established")
	return client, nil
}
