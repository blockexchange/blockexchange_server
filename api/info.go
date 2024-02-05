package api

import (
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
	return nil
}

func (api *Api) GetInfo(w http.ResponseWriter, req *http.Request) {

	cfg := api.cfg
	info := Info{
		VersionMajor: 1,
		VersionMinor: 1,
		Name:         cfg.Name,
		Owner:        cfg.Owner,
		BaseURL:      cfg.BaseURL,
		OauthLogin:   cfg.OauthLogin,
		Stats:        stats.Load(),
	}

	SendJson(w, info)
}
