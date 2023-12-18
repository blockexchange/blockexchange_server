package core

import (
	"blockexchange/db"
	"blockexchange/types"
)

type Core struct {
	cfg   *types.Config
	repos *db.Repositories
}

func New(cfg *types.Config, repos *db.Repositories) *Core {
	return &Core{
		cfg:   cfg,
		repos: repos,
	}
}
