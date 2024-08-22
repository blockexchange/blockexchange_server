package jobs

import (
	"blockexchange/api"
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/types"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	LOCKID_CLEANUP      = 100
	LOCKID_STATS_UPDATE = 101
	LOCKID_SCREENSHOT   = 102
)

func Start(repos *db.Repositories, api *api.Api) {
	cfg := types.CreateConfig()
	c := core.New(cfg, repos)

	if cfg.ExecuteJobs {
		// start jobs
		go cleanupSchemas(repos.SchemaRepo, repos.Lock)
		go updateScreenshots(c, api.SchemaSearchRepo, repos.Lock)
	}
	// per-node stats job
	go updateStats(api, repos.Lock)

}

func updateStats(api *api.Api, lock *db.DBLock) {
	for {
		err := lock.RunLocked(LOCKID_STATS_UPDATE, time.Minute, api.UpdateStats)
		if err != nil {
			logrus.WithError(err).Error("stats update")
		}
		time.Sleep(30 * time.Minute)
	}
}

func cleanupSchemas(schemarepo *db.SchemaRepository, lock *db.DBLock) {
	for {
		err := lock.RunLocked(LOCKID_CLEANUP, time.Minute, func() error {
			logrus.Trace("Removing old and incomplete schemas")
			now := time.Now().Unix() * 1000
			return schemarepo.DeleteOldIncompleteSchema(now - (3600 * 1000 * 24))
		})
		if err != nil {
			logrus.WithError(err).Error("cleanupSchemas")
		}
		time.Sleep(5 * time.Minute)
	}
}

func updateScreenshots(c *core.Core, sr *db.SchemaSearchRepository, lock *db.DBLock) {
	from := time.Now().Add(-10*time.Minute).Unix() * 1000

	for {
		logrus.Trace("updating schema previews")
		err := lock.RunLocked(LOCKID_SCREENSHOT, time.Minute, func() error {
			complete := true
			list, err := sr.Search(&types.SchemaSearchRequest{
				FromMtime: &from,
				Complete:  &complete,
			})
			if err != nil {
				return fmt.Errorf("schema search error: %v", err)
			}

			for _, r := range list {
				logrus.WithFields(logrus.Fields{
					"uid":   r.Schema.UID,
					"mtime": r.Schema.Mtime,
				}).Debug("Updating schema screenshot")
				_, err = c.UpdatePreview(r.Schema)
				if err != nil {
					return fmt.Errorf("schema preview update error: '%s', %v", r.Schema.UID, err)
				}

				// shift mtime window to max mtime from result
				if r.Schema.Mtime > from {
					from = r.Schema.Mtime
				}
			}
			return nil
		})
		if err != nil {
			logrus.WithError(err).Errorf("schema screenshot update: %d", from)
			continue
		}

		time.Sleep(5 * time.Minute)
	}
}
