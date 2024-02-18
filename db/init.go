package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5"
	"github.com/vingarcia/ksql"
	"github.com/vingarcia/ksql/adapters/kpgx"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
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

	ctx := context.Background()
	kdb, err := kpgx.New(ctx, url, ksql.Config{})
	if err != nil {
		return nil, err
	}

	return NewRepositories(kdb, db), nil
}
