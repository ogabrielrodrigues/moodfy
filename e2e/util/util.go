package util

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func TestDatabase() *sql.DB {
	db_url := os.Getenv("TEST_DB_URL")

	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

func ClearDatabase(db *sql.DB, tb string) {
	db.Exec(fmt.Sprintf("DELETE FROM %s", tb))
}
