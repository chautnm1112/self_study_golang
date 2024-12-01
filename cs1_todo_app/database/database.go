package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
)

var db *sql.DB

// Global logger for the application
var logger *zap.Logger

func GetDB() *sql.DB {
	return db
}

// InitDB initializes the database connection and sets up the logger
func InitDB(connStr string) {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		logger.Fatal("Error initializing zap logger: %v", zap.Error(err))
	}

	// Open a connection to the database
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatal("Can't connect to database", zap.Error(err))
	}

	// Verify the database connection
	err = db.Ping()
	if err != nil {
		logger.Fatal("Database ping error", zap.Error(err))
	}

	// Log success message
	logger.Info("Database connected successfully!")

	_, err = db.Exec("DROP TABLE IF EXISTS task")
	if err != nil {
		log.Fatalf("Failed to drop task table: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE task (
		task_id SERIAL PRIMARY KEY,
		task TEXT,
		completed BOOLEAN,
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	)`)
	
	if err != nil {
		log.Fatalf("Failed to create task table: %v", err)
	}
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			logger.Error("Error closing DB connection", zap.Error(err))
		} else {
			logger.Info("Database connection closed!")
		}
	}
}
