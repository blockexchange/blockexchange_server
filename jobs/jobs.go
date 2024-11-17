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

func updateStats(api *api.Api, pgl *pglock.Client) {
	for {
		lock, err := pgl.Acquire("update-stats")
		if err != nil {
			logrus.WithError(err).Error("update stats lock")
			time.Sleep(time.Second * 10)
			continue
		}

		err = api.UpdateStats()
		if err != nil {
			logrus.WithError(err).Error("update stats")
		}

		lock.Close()
		time.Sleep(30 * time.Minute)
	}
}

func cleanupSchemas(schemarepo *db.SchemaRepository, pgl *pglock.Client) {
	for {
		lock, err := pgl.Acquire("schema-cleanup")
		if err != nil {
			logrus.WithError(err).Error("schema cleanup lock")
			time.Sleep(time.Second * 10)
			continue
		}

		logrus.Trace("Removing old and incomplete schemas")
		now := time.Now().Unix() * 1000
		err = schemarepo.DeleteOldIncompleteSchema(now - (3600 * 1000 * 24))
		if err != nil {
			logrus.WithError(err).Error("schema cleanup")
		}

		lock.Close()
		time.Sleep(5 * time.Minute)
	}
}

func updateScreenshots(c *core.Core, sr *db.SchemaSearchRepository, pgl *pglock.Client) {
	from := time.Now().Add(-10*time.Minute).Unix() * 1000

	for {
		lock, err := pgl.Acquire("update-screenshots")
		if err != nil {
			logrus.WithError(err).Error("screenshot update lock")
			time.Sleep(time.Second * 10)
			continue
		}

		logrus.Trace("updating schema previews")
		complete := true
		list, err := sr.Search(&types.SchemaSearchRequest{
			FromMtime: &from,
			Complete:  &complete,
		})
		if err != nil {
			logrus.WithError(err).Error("schema search")
			time.Sleep(time.Second * 10)
			continue
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

		lock.Close()
		time.Sleep(5 * time.Minute)
	}
}
