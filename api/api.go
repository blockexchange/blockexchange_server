package api

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/types"
	"sync/atomic"

	"github.com/minetest-go/colormapping"

	"github.com/jmoiron/sqlx"
)

type Api struct {
	*db.Repositories
	cfg          *types.Config
	core         *core.Core
	Cache        core.Cache
	ColorMapping *colormapping.ColorMapping
	Running      *atomic.Bool
}

func (a *Api) Stop() {
	a.Running.Store(false)
}

func NewApi(db_ *sqlx.DB, cache core.Cache, cfg *types.Config) (*Api, error) {
	cm := colormapping.NewColorMapping()
	err := cm.LoadDefaults()
	if err != nil {
		return nil, err
	}

	running := &atomic.Bool{}
	running.Store(true)

	repos := db.NewRepositories(db_)

	return &Api{
		Repositories: repos,
		cfg:          cfg,
		core:         core.New(cfg, repos),
		Cache:        cache,
		ColorMapping: cm,
		Running:      running,
	}, nil
}
