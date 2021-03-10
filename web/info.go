package web

import (
	"encoding/json"
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
}

func InfoEndpoint(w http.ResponseWriter, req *http.Request) {
	info := Info{
		VersionMajor: 1,
		VersionMinor: 1,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(info)
}
