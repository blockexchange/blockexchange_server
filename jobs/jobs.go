package jobs

import (
	"blockexchange/api"
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/types"
	"time"

	"github.com/sirupsen/logrus"
)

func Start(repos *db.Repositories, api *api.Api) {
	cfg := types.CreateConfig()
	c := core.New(cfg, repos)

	if cfg.ExecuteJobs {
		// start jobs
		go cleanupSchemas(repos.SchemaRepo)
		go updateScreenshots(c, api.SchemaSearchRepo)
	}
	// per-node stats job
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

func cleanupSchemas(schemarepo *db.SchemaRepository) {
	for {
		logrus.Trace("Removing old and incomplete schemas")
		now := time.Now().Unix() * 1000
		schemarepo.DeleteOldIncompleteSchema(now - (3600 * 1000 * 24))
		time.Sleep(5 * time.Minute)
	}
}

func updateScreenshots(c *core.Core, sr *db.SchemaSearchRepository) {
	from := time.Now().Add(-10*time.Minute).Unix() * 1000

	for {
		logrus.Trace("updating schema previews")
		complete := true
		list, err := sr.Search(&types.SchemaSearchRequest{
			FromMtime: &from,
			Complete:  &complete,
		})
		if err != nil {
			logrus.WithError(err).Errorf("schema mtime search from: %d", from)
			continue
		}

		for _, r := range list {
			logrus.WithFields(logrus.Fields{
				"uid":   r.Schema.UID,
				"mtime": r.Schema.Mtime,
			}).Debug("Updating schema screenshot")
			_, err = c.UpdatePreview(r.Schema)
			if err != nil {
				logrus.WithError(err).Errorf("schema preview update: '%s'", r.Schema.UID)
				continue
			}

			// shift mtime window to max mtime from result
			if r.Schema.Mtime > from {
				from = r.Schema.Mtime
			}
		}

		time.Sleep(5 * time.Minute)
	}
}
