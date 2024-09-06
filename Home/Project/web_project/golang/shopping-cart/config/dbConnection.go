package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	var err error
	DB, err = sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Failed to open the database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	fmt.Println("Successfully connected to database")
}

func CreateTables() {
	const createTablesSQL = `
	CREATE TABLE IF NOT EXISTS cart_items (
		item_id INT PRIMARY KEY,
		item_name TEXT NOT NULL,
		price FLOAT NOT NULL,
		quantity INT NOT NULL
	);
	`
	if _, err := DB.Exec(createTablesSQL); err != nil {
		log.Fatal("Failed to create tables:", err)
	}
}
