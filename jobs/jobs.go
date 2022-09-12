package jobs

import (
	"blockexchange/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func Start(db_ *sqlx.DB) {
	schemarepo := db.DBSchemaRepository{DB: db_}
	go loop(schemarepo)
}

func loop(schemarepo db.DBSchemaRepository) {
	for {
		cleanupSchemas(schemarepo)
		time.Sleep(5 * time.Minute)
	}
}

func cleanupSchemas(schemarepo db.DBSchemaRepository) {
	logrus.Trace("Removing old and incomplete schemas")
	now := time.Now().Unix() * 1000
	schemarepo.DeleteOldIncompleteSchema(now - (3600 * 1000 * 24))
}
