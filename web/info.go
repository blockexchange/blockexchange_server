package web

import (
	"net/http"
	"os"
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
	oauth := OauthInfo{
		GithubID:  os.Getenv("GITHUB_APP_ID"),
		DiscordID: os.Getenv("DISCORD_APP_ID"),
		MesehubID: os.Getenv("MESEHUB_APP_ID"),
		BaseURL:   os.Getenv("BASE_URL"),
	}

	info := Info{
		VersionMajor: 1,
		VersionMinor: 1,
		Name:         os.Getenv("BLOCKEXCHANGE_NAME"),
		Owner:        os.Getenv("BLOCKEXCHANGE_OWNER"),
		Oauth:        &oauth,
	}

	SendJson(w, info)
}
