package oauth

import (
	"blockexchange/core"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type DiscordResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type DiscordOauth struct {
}

func (o *DiscordOauth) RequestAccessToken(code string, cfg *core.Config) (string, error) {
	q := url.Values{}
	q.Add("client_id", cfg.DiscordOAuthConfig.ClientID)
	q.Add("client_secret", cfg.DiscordOAuthConfig.Secret)
	q.Add("redirect_uri", cfg.BaseURL+"/api/oauth_callback/discord")
	q.Add("code", code)
	q.Add("grant_type", "authorization_code")
	q.Add("scope", "identify email connections")
	fmt.Println(q.Encode())

	buf := bytes.NewBufferString(q.Encode())

	req, err := http.NewRequest("POST", "https://discord.com/api/oauth2/token", buf)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	tokenData := AccessTokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&tokenData)
	if err != nil {
		return "", err
	}

	return tokenData.AccessToken, nil
}

func (o *DiscordOauth) RequestUserInfo(access_token string, cfg *core.Config) (*OauthUserInfo, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		return nil, nil
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+access_token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	userData := DiscordResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		return nil, err
	}

	fmt.Println(userData)
	info := OauthUserInfo{
		Name:       userData.Username,
		Type:       types.UserTypeDiscord,
		Email:      userData.Email,
		ExternalID: userData.ID,
	}

	return &info, nil
}
