package jobs

import (
	"blockexchange/api"
	"blockexchange/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func Start(db_ *sqlx.DB, api *api.Api) {
	schemarepo := db.SchemaRepository{DB: db_.DB}
	go cleanupSchemas(schemarepo)
	go updateStats(api)
}

func updateStats(api *api.Api) {
	for {
		err := api.UpdateStats()
		if err != nil {
			logrus.WithError(err).Error("stats update")
		}
		time.Sleep(30 * time.Minute)
	}
}

func cleanupSchemas(schemarepo db.SchemaRepository) {
	for {
		logrus.Trace("Removing old and incomplete schemas")
		now := time.Now().Unix() * 1000
		schemarepo.DeleteOldIncompleteSchema(now - (3600 * 1000 * 24))
		time.Sleep(5 * time.Minute)
	}
}
