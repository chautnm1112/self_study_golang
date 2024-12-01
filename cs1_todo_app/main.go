package main

import (
	"cs1_todo_app/database"
	"cs1_todo_app/routes"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"net/http"
	"os"
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
	// Start the server
	err = http.ListenAndServe(":"+host, r)
	if err != nil {
		// Log the error with zap
		logger.Fatal("Can't run server", zap.Error(err))
	}
}
