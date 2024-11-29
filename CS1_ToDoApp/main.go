package main

import (
	"CS1_ToDoApp/database"
	"CS1_ToDoApp/routes"
	"go.uber.org/zap"
	"net/http"
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
	defer logger.Sync() // Ensure buffered logs are flushed

	// Initialize the database
	database.InitDB()
	defer database.CloseDB()

	// Set up router
	r := routes.SetupRouter()

	// Start the server
	err = http.ListenAndServe(":8888", r)
	if err != nil {
		// Log the error with zap
		logger.Fatal("Can't run server", zap.Error(err))
	}
}
