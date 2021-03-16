package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Init() (*sqlx.DB, error) {
	connStr := fmt.Sprintf(
		"user=%s password=%s port=%s host=%s dbname=%s sslmode=disable",
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGPORT"),
		os.Getenv("PGHOST"),
		os.Getenv("PGDATABASE"))

	log.Printf("Connecting to %s", connStr)
	var err error
	DB, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = DB.Ping()
	if err != nil {
		return nil, err
	}

	return DB, nil
}
