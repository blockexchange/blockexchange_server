package core

import (
	"blockexchange/db"
	"blockexchange/types"
)

type Core struct {
	cfg       *types.Config
	repos     *db.Repositories
	wasmPacks *wasmPackCache
}

func New(cfg *types.Config, repos *db.Repositories) *Core {
	return &Core{
		cfg:       cfg,
		repos:     repos,
		wasmPacks: newWasmPackCache(cfg.WasmPackCacheDir, cfg.WasmGitCacheDir),
	}
}
