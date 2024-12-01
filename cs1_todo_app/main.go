package main

import (
	"context"
	"cs1_todo_app/database"
	"cs1_todo_app/routes"
	"errors"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var logger *zap.Logger

func main() {
	// Initialize zap logger
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		// If logger can't be initialized, log to standard output and stop
		logger.Fatal("Error initializing zap logger: %v", zap.Error(err))
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {

		}
	}(logger) // Ensure buffered logs are flushed

	err = godotenv.Load(".env")
	if err != nil {
		logger.Fatal("Error loading .env file: ", zap.Error(err))
	}

	dsn := os.Getenv("APP_ENV")
	// Initialize the database
	database.InitDB(dsn)
	defer database.CloseDB()

	// Set up router
	r := routes.SetupRouter()

	host := os.Getenv("HOST")

	server := &http.Server{
		Addr:    ":" + host,
		Handler: r,
	}

	timeWait := 15 * time.Second

	signChan := make(chan os.Signal, 1)

	go func() {
		logger.Info("Starting server on :" + host)
		// Start the server
		if err := server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logger.Error("Could not listen on "+server.Addr+": ", zap.Error(err))
		}
	}()

	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
	<-signChan

	logger.Info("Shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), timeWait)
	defer func() {
		logger.Info("Close another connection")
		cancel()
	}()

	logger.Info("Stop http server")

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", zap.Error(err))
	}

	close(signChan)
	logger.Info("Server stopped gracefully")
}
