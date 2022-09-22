package jobs

import (
	"blockexchange/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func Start(db_ *sqlx.DB) {
	schemarepo := db.SchemaRepository{DB: db_.DB}
	go loop(schemarepo)
}

func loop(schemarepo db.SchemaRepository) {
	for {
		cleanupSchemas(schemarepo)
		time.Sleep(5 * time.Minute)
	}
}

func cleanupSchemas(schemarepo db.SchemaRepository) {
	logrus.Trace("Removing old and incomplete schemas")
	now := time.Now().Unix() * 1000
	schemarepo.DeleteOldIncompleteSchema(now - (3600 * 1000 * 24))
}
