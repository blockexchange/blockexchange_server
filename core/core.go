package core

import "blockexchange/types"

type Core struct {
	cfg *types.Config
}

func New(cfg *types.Config) *Core {
	return &Core{
		cfg: cfg,
	}
}
