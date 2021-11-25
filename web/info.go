package web

import (
	"blockexchange/core"
	"net/http"
)

type OauthInfo struct {
	GithubID  string `json:"github_id"`
	DiscordID string `json:"discord_id"`
	MesehubID string `json:"mesehub_id"`
	BaseURL   string `json:"base_url"`
}

type Info struct {
	VersionMajor int        `json:"api_version_major"`
	VersionMinor int        `json:"api_version_minor"`
	Name         string     `json:"name"`
	Owner        string     `json:"owner"`
	Oauth        *OauthInfo `json:"oauth"`
	EnableSignup bool       `json:"enable_signup"`
}

type InfoHandler struct {
	Config *core.Config
}

func (h InfoHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	oauth := OauthInfo{
		BaseURL: h.Config.BaseURL,
	}

	if h.Config.GithubOAuthConfig != nil {
		oauth.GithubID = h.Config.GithubOAuthConfig.ClientID
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
		Oauth:        &oauth,
		EnableSignup: h.Config.EnableSignup,
	}

	SendJson(w, info)
}
