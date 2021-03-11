package db

import (
	"embed"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrations embed.FS

func init() {
	source.Register("embed", &driver{})
}

type driver struct {
	httpfs.PartialDriver
}

func (d *driver) Open(rawURL string) (source.Driver, error) {
	err := d.PartialDriver.Init(http.FS(migrations), "migrations")
	if err != nil {
		return nil, err
	}

	return d, nil
}

func Migrate() {
	driver, err := postgres.WithInstance(DB.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance("embed://", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	v, _, _ := m.Version()
	log.Printf("DB-Version: %d", v)
}
