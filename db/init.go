package db

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"time"

	"cirello.io/pglock"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	mpg "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Init() (*Repositories, error) {
	// main pg url
	url := fmt.Sprintf(
		"user=%s password=%s port=%s host=%s dbname=%s sslmode=disable",
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGPORT"),
		os.Getenv("PGHOST"),
		os.Getenv("PGDATABASE"))

	fmt.Printf("Connecting to %s\n", url)

	// gorm instance

	g, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm open: %v", err)
	}

	// db migrations

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	driver, err := mpg.WithInstance(db, &mpg.Config{})
	if err != nil {
		return nil, err
	}

	d, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance("iofs", d, "pgx", driver)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// pglock

	pgl, err := pglock.New(db, pglock.WithLeaseDuration(3*time.Second), pglock.WithHeartbeatFrequency(1*time.Second))
	if err != nil {
		return nil, fmt.Errorf("cannot create lock client: %v", err)
	}

	err = pgl.TryCreateTable()
	if err != nil {
		return nil, fmt.Errorf("cannot create table: %v", err)
	}

	return NewRepositories(g, db, pgl), nil
}
