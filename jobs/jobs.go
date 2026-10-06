package jobs

import (
	"blockexchange/api"
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/types"
	"time"

	"cirello.io/pglock"
	"github.com/sirupsen/logrus"
)

func Start(repos *db.Repositories, api *api.Api) {
	cfg := types.CreateConfig()
	c := core.New(cfg, repos)

	if cfg.ExecuteJobs {
		// start jobs
		go cleanupSchemas(repos.SchemaRepo, repos.PGLock)
		go updateScreenshots(c, api.SchemaSearchRepo, repos.PGLock)
	}
	// per-node stats job
	go updateStats(api, repos.PGLock)

}

// runLocked executes fn every interval while holding the named cluster-wide lock
func runLocked(pgl *pglock.Client, name string, interval time.Duration, fn func()) {
	for {
		lock, err := pgl.Acquire(name)
		if err != nil {
			logrus.WithError(err).WithField("job", name).Error("job lock")
			time.Sleep(time.Second * 10)
			continue
		}

		fn()

		lock.Close()
		time.Sleep(interval)
	}
}

func updateStats(api *api.Api, pgl *pglock.Client) {
	runLocked(pgl, "update-stats", 30*time.Minute, func() {
		err := api.UpdateStats()
		if err != nil {
			logrus.WithError(err).Error("update stats")
		}
	})
}

func cleanupSchemas(schemarepo *db.SchemaRepository, pgl *pglock.Client) {
	runLocked(pgl, "schema-cleanup", 5*time.Minute, func() {
		logrus.Trace("Removing old and incomplete schemas")
		now := time.Now().Unix() * 1000
		err := schemarepo.DeleteOldIncompleteSchema(now - (3600 * 1000 * 24))
		if err != nil {
			logrus.WithError(err).Error("schema cleanup")
		}
	})
}

func updateScreenshots(c *core.Core, sr *db.SchemaSearchRepository, pgl *pglock.Client) {
	from := time.Now().Add(-10*time.Minute).Unix() * 1000

	runLocked(pgl, "update-screenshots", 5*time.Minute, func() {
		logrus.Trace("updating schema previews")
		complete := true
		list, err := sr.Search(&types.SchemaSearchRequest{
			FromMtime: &from,
			Complete:  &complete,
		})
		if err != nil {
			logrus.WithError(err).Error("schema search")
			return
		}

		for _, r := range list {
			logrus.WithFields(logrus.Fields{
				"uid":   r.Schema.UID,
				"mtime": r.Schema.Mtime,
			}).Debug("Updating schema screenshot")
			_, err = c.UpdatePreview(r.Schema)
			if err != nil {
				logrus.Errorf("schema preview update error: '%s', %v", r.Schema.UID, err)
				continue
			}

			// shift mtime window to max mtime from result
			if r.Schema.Mtime > from {
				from = r.Schema.Mtime
			}
		}
	})
}
