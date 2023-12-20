package oauth

import (
	"net/http"
)

type OauthHandler struct {
	provider OauthProvider
	cfg      *OAuthConfig
	cb       OauthCallback
}

func NewHandler(cb OauthCallback, cfg *OAuthConfig) *OauthHandler {
	var provider OauthProvider
	switch cfg.Provider {
	case ProviderTypeCDB:
		provider = &CDBOauth{}
	case ProviderTypeDiscord:
		provider = &DiscordOauth{}
	case ProviderTypeGithub:
		provider = &GithubOauth{}
	case ProviderTypeMesehub:
		provider = &MesehubOauth{}
	default:
		panic("unkown provider-type: " + cfg.Provider)
	}

	return &OauthHandler{provider: provider, cfg: cfg, cb: cb}
}

func (h *OauthHandler) LoginURL() string {
	return h.provider.LoginURL(h.cfg)
}

func (h *OauthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list := r.URL.Query()["code"]
	if len(list) == 0 {
		SendError(w, 500, "no code found")
		return
	}

	code := list[0]

	access_token, err := h.provider.RequestAccessToken(code, h.cfg)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	info, err := h.provider.RequestUserInfo(access_token, h.cfg)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if info.ExternalID == "" {
		SendError(w, 500, "empty external_id")
		return
	}

	err = h.cb(w, r, info)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

}
