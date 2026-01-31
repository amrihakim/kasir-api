package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Initialize your database connection here
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	// Test the connection
	// err = db.Ping()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Set connection pool settings if needed
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	fmt.Println("Database connection established")
	return db, nil
}
