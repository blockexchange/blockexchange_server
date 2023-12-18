package api

import (
	"blockexchange/types"
	"net/http"
)

type OauthInfo struct {
	GithubID  string `json:"github_id"`
	DiscordID string `json:"discord_id"`
	MesehubID string `json:"mesehub_id"`
	CDBID     string `json:"cdb_id"`
}

type Info struct {
	VersionMajor int        `json:"api_version_major"`
	VersionMinor int        `json:"api_version_minor"`
	Name         string     `json:"name"`
	Owner        string     `json:"owner"`
	BaseURL      string     `json:"base_url"`
	Oauth        *OauthInfo `json:"oauth"`
}

type InfoHandler struct {
	Config *types.Config
}

func (h InfoHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	oauth := OauthInfo{}

	if h.Config.GithubOAuthConfig != nil {
		oauth.GithubID = h.Config.GithubOAuthConfig.ClientID
	}

	if h.Config.CDBOAuthConfig != nil {
		oauth.CDBID = h.Config.CDBOAuthConfig.ClientID
	}

	if h.Config.DiscordOAuthConfig != nil {
		oauth.DiscordID = h.Config.DiscordOAuthConfig.ClientID
	}

	if h.Config.MesehubOAuthConfig != nil {
		oauth.MesehubID = h.Config.MesehubOAuthConfig.ClientID
	}

	info := Info{
		VersionMajor: 1,
		VersionMinor: 1,
		Name:         h.Config.Name,
		Owner:        h.Config.Owner,
		BaseURL:      h.Config.BaseURL,
		Oauth:        &oauth,
	}

	SendJson(w, info)
}
