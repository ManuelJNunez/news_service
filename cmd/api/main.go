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
	"github.com/ManuelJNunez/news_service/internal/news"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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

	// 3) Connect to database
	db, err := initDB(cfg, logger)

	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	// 4) Build dependencies from news domain
	newsRepo := news.NewRepository(db)
	newsSvc := news.NewService(newsRepo)
	newsHandler := news.NewHandler(newsSvc)

	// 5) Configure Gin (web framework)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router_group := router.Group("")

	news.RegisterRoutes(router_group, newsHandler)

	// 6) Configure HTTP server
	addr := ":" + cfg.HTTPPort
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// 7) Launch webserver on a separate Thread using a Goroutine
	go func() {
		logger.Info("HTTP server listening", slog.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server listen error", slog.Any("error", err))
		}
	}()

	// 8) Configure signal handling and pause execution until SIGINT or SIGTERM is received
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
	// Open DB connection (PostgreSQL)
	db, err := sql.Open("postgres", cfg.DB_DSN)

	if err != nil {
		return nil, err
	}

	// Initialize context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping database connection, to check if connection is alive
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		logger.Error("database ping failed", slog.Any("error", err))
		return nil, err
	}

	logger.Info("database connection established")
	return db, nil
}
