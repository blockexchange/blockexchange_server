package api

import (
	"blockexchange/colormapping"
	"blockexchange/core"
	"blockexchange/db"

	"github.com/jmoiron/sqlx"
)

type Api struct {
	*db.Repositories
	Cache        core.Cache
	ColorMapping *colormapping.ColorMapping
}

func NewApi(db_ *sqlx.DB, cache core.Cache) (*Api, error) {
	cm := colormapping.NewColorMapping()
	err := cm.LoadDefaults()
	if err != nil {
		return nil, err
	}

	return &Api{
		Repositories: db.NewRepositories(db_),
		Cache:        cache,
		ColorMapping: cm,
	}, nil
}
