package api

import (
	"blockexchange/types"
	"net/http"
)

type Info struct {
	VersionMajor int               `json:"api_version_major"`
	VersionMinor int               `json:"api_version_minor"`
	Name         string            `json:"name"`
	Owner        string            `json:"owner"`
	BaseURL      string            `json:"base_url"`
	OauthLogin   *types.OauthLogin `json:"oauth_login"`
}

type InfoHandler struct {
	Config *types.Config
}

func (h InfoHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	info := Info{
		VersionMajor: 1,
		VersionMinor: 1,
		Name:         h.Config.Name,
		Owner:        h.Config.Owner,
		BaseURL:      h.Config.BaseURL,
		OauthLogin:   h.Config.OauthLogin,
	}

	SendJson(w, info)
}
