package db

import (
	"database/sql"
	"embed"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Init() (*Repositories, error) {
	url := fmt.Sprintf(
		"user=%s password=%s port=%s host=%s dbname=%s sslmode=disable",
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGPORT"),
		os.Getenv("PGHOST"),
		os.Getenv("PGDATABASE"))

	fmt.Printf("Connecting to %s\n", url)

	g, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm open: %v", err)
	}

	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	driver, err := pgx.WithInstance(db, &pgx.Config{})
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

	return NewRepositories(g, db), nil
}
