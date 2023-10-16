package repository

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func SetupDatabaseSchema(db *sql.DB) {
	schema := `
    CREATE TABLE IF NOT EXISTS movies (
        id INT AUTO_INCREMENT PRIMARY KEY,
		titleType VARCHAR(50) NOT NULL,
        title VARCHAR(255) NOT NULL,
        year YEAR NOT NULL,
        genres VARCHAR(200),
        runtime INT,
		rating DECIMAL(1,1),
		votes INT
    );`

	_, err := db.Exec(schema)
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}
}

func NewDatabase() *sql.DB {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dsn := os.Getenv("DSN") // Fetch the DSN from environment variables
	if dsn == "" {
		log.Fatal("DSN not set in .env")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	return db
}
