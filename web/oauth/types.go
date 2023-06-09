package oauth

import (
	"blockexchange/core"
	"blockexchange/types"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type OauthUserInfo struct {
	Name       string
	Email      string
	ExternalID string
	Type       types.UserType
}

type OauthImplementation interface {
	RequestAccessToken(code string, cfg *core.Config) (string, error)
	RequestUserInfo(access_token string, cfg *core.Config) (*OauthUserInfo, error)
}
