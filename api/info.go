package api

import (
	"blockexchange/metrics"
	"blockexchange/types"
	"net/http"
	"sync/atomic"
)

type Info struct {
	VersionMajor int               `json:"api_version_major"`
	VersionMinor int               `json:"api_version_minor"`
	Name         string            `json:"name"`
	Owner        string            `json:"owner"`
	BaseURL      string            `json:"base_url"`
	OauthLogin   *types.OauthLogin `json:"oauth_login"`
	Stats        *Stats            `json:"stats"`
	// luanti wasm client available ("Edit Online")
	WasmEnabled bool `json:"wasm_enabled"`
}

type Stats struct {
	SchemaCount     int64 `json:"schema_count"`
	SchemaPartCount int64 `json:"schemapart_count"`
	UserCount       int64 `json:"user_count"`
	TotalSize       int64 `json:"total_size"`
}

type InfoHandler struct {
	Config *types.Config
}

var stats atomic.Pointer[Stats]

func init() {
	stats.Store(&Stats{})
}

func (api *Api) UpdateStats() error {
	new_stats := &Stats{}
	var err error

	new_stats.SchemaCount, err = api.MetaRepository.CountEntries("schema")
	if err != nil {
		return err
	}

	new_stats.SchemaPartCount, err = api.MetaRepository.CountEntries("schemapart")
	if err != nil {
		return err
	}

	new_stats.UserCount, err = api.MetaRepository.CountEntries("user")
	if err != nil {
		return err
	}

	new_stats.TotalSize, err = api.SchemaRepo.GetTotalSize()
	if err != nil {
		return err
	}

	stats.Store(new_stats)

	metrics.Schemas.Set(float64(new_stats.SchemaCount))
	metrics.SchemaParts.Set(float64(new_stats.SchemaPartCount))
	metrics.Users.Set(float64(new_stats.UserCount))
	metrics.StorageBytes.Set(float64(new_stats.TotalSize))
	return nil
}

func (api *Api) GetInfo(w http.ResponseWriter, req *http.Request) {

	cfg := api.cfg
	info := Info{
		VersionMajor: 1,
		VersionMinor: 2,
		Name:         cfg.Name,
		Owner:        cfg.Owner,
		BaseURL:      cfg.BaseURL,
		OauthLogin:   cfg.OauthLogin,
		Stats:        stats.Load(),
		WasmEnabled:  true,
	}

	SendJson(w, info)
}
