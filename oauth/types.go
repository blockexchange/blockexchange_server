package oauth

import (
	"net/http"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type OauthUserInfo struct {
	Provider    ProviderType
	Name        string
	ExternalID  string
	AvatarURL   string
	DisplayName string
}

type OAuthConfig struct {
	Provider    ProviderType
	ClientID    string
	Secret      string
	CallbackURL string
}

type ProviderType string

const (
	ProviderTypeGithub  ProviderType = "GITHUB"
	ProviderTypeDiscord ProviderType = "DISCORD"
	ProviderTypeMesehub ProviderType = "MESEHUB"
	ProviderTypeCDB     ProviderType = "CDB"
)

type OauthProvider interface {
	// client login url
	LoginURL(cfg *OAuthConfig) string
	RequestAccessToken(code string, cfg *OAuthConfig) (string, error)
	RequestUserInfo(access_token string, cfg *OAuthConfig) (*OauthUserInfo, error)
}

type OauthCallback func(w http.ResponseWriter, r *http.Request, user_info *OauthUserInfo) error
