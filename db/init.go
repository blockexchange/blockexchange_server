package db

import (
	"context"
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

// connection pool limits
var (
	MaxOpenConns = 20
	MaxIdleConns = 5
)

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

	// one pool for gorm, the migrations, the plain sql queries and pglock (lib/pq: pglock and pq.Array depend on it)
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetConnMaxIdleTime(5 * time.Minute)

	err = migrateUp(db)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	g, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm open: %v", err)
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

// runs the migrations on a dedicated connection, released to the pool afterwards
func migrateUp(db *sql.DB) error {
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	// not driver.Close(): it would close the shared pool
	defer conn.Close()

	driver, err := mpg.WithConnection(ctx, conn, &mpg.Config{})
	if err != nil {
		return err
	}

	d, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
