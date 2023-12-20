package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type DiscordResponse struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	AvatarHash string `json:"avatar"`
	GlobalName string `json:"global_name"`
}

type DiscordOauth struct{}

func (o *DiscordOauth) LoginURL(cfg *OAuthConfig) string {
	return fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify", cfg.ClientID, url.QueryEscape(cfg.CallbackURL))
}

func (o *DiscordOauth) RequestAccessToken(code string, cfg *OAuthConfig) (string, error) {
	q := url.Values{}
	q.Add("client_id", cfg.ClientID)
	q.Add("client_secret", cfg.Secret)
	q.Add("redirect_uri", cfg.CallbackURL)
	q.Add("code", code)
	q.Add("grant_type", "authorization_code")
	q.Add("scope", "identify connections")

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
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("invalid status code in token-response: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	tokenData := AccessTokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&tokenData)
	if err != nil {
		return "", err
	}

	return tokenData.AccessToken, nil
}

func (o *DiscordOauth) RequestUserInfo(access_token string, cfg *OAuthConfig) (*OauthUserInfo, error) {
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
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code in response: %d", resp.StatusCode)
	}

	userData := DiscordResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		return nil, err
	}

	info := OauthUserInfo{
		Provider:    ProviderTypeDiscord,
		Name:        userData.Username,
		ExternalID:  userData.ID,
		DisplayName: userData.GlobalName,
		AvatarURL:   fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", userData.ID, userData.AvatarHash),
	}

	return &info, nil
}
