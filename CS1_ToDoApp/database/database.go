package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var Db *sql.DB

func InitDB() {
	connStr := "user=postgres dbname=todoApp password=123321 host=localhost sslmode=disable"
	var err error
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Cann't connect to database", err)
	}

	err = Db.Ping()
	if err != nil {
		log.Fatal("Database ping error: ", err)
	}

	log.Println("Database connected successfully!")
}

func CloseDB() {
	if Db != nil {
		if err := Db.Close(); err != nil {
			log.Printf("Error closing DB connection: %v", err)
		} else {
			log.Println("Database connection closed!")
		}
	}
}
