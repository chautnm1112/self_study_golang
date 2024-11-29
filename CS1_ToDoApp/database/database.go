package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
)

var Db *sql.DB

// Global logger for the application
var logger *zap.Logger

// InitDB initializes the database connection and sets up the logger
func InitDB() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("Error initializing zap logger: %v", err)
	}

	// Set up the connection string for PostgreSQL
	connStr := "user=postgres dbname=todoApp password=123321 host=localhost sslmode=disable"

	// Open a connection to the database
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatal("Can't connect to database", zap.Error(err))
	}

	// Verify the database connection
	err = Db.Ping()
	if err != nil {
		logger.Fatal("Database ping error", zap.Error(err))
	}

	// Log success message
	logger.Info("Database connected successfully!")
}

// CloseDB closes the database connection
func CloseDB() {
	if Db != nil {
		if err := Db.Close(); err != nil {
			logger.Error("Error closing DB connection", zap.Error(err))
		} else {
			logger.Info("Database connection closed!")
		}
	}
}
