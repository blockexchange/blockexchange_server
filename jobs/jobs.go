package jobs

import (
	"blockexchange/core"
	"blockexchange/db"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func Start(db_ *sqlx.DB) {
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")

	var syncmanager core.SyncManager
	if redis_host != "" && redis_port != "" {
		syncmanager = core.NewRedisSyncManager(redis_host + ":" + redis_port)
	} else {
		syncmanager = core.NewLocalSyncManager()
	}

	schemarepo := db.DBSchemaRepository{DB: db_}

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(5).Minutes().Do(cleanupSchemas, schemarepo, syncmanager)
}

func cleanupSchemas(schemarepo db.DBSchemaRepository, syncmanager core.SyncManager) {
	lock := syncmanager.GetLock("schema-cleanup")
	lock.Lock()
	defer lock.Unlock()

	logrus.Debugf("Removing old schemas")
	now := time.Now().Unix() * 1000
	schemarepo.DeleteOldIncompleteSchema(now - (3600 * 1000 * 24))
}
