package web

import (
	"blockexchange/core"
	"blockexchange/db"

	"github.com/jmoiron/sqlx"
)

type Api struct {
	*db.Repositories
	Cache core.Cache
}

func NewApi(db_ *sqlx.DB, cache core.Cache) *Api {
	return &Api{
		Repositories: db.NewRepositories(db_),
		Cache:        cache,
	}
}
