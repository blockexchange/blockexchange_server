package api

import (
	"blockexchange/core"
	"blockexchange/db"
	"sync/atomic"

	"github.com/minetest-go/colormapping"

	"github.com/jmoiron/sqlx"
)

type Api struct {
	*db.Repositories
	Cache        core.Cache
	ColorMapping *colormapping.ColorMapping
	Running      *atomic.Bool
}

func (a *Api) Stop() {
	a.Running.Store(false)
}

func NewApi(db_ *sqlx.DB, cache core.Cache) (*Api, error) {
	cm := colormapping.NewColorMapping()
	err := cm.LoadDefaults()
	if err != nil {
		return nil, err
	}

	running := &atomic.Bool{}
	running.Store(true)

	return &Api{
		Repositories: db.NewRepositories(db_),
		Cache:        cache,
		ColorMapping: cm,
		Running:      running,
	}, nil
}
