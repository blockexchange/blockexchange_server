package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type GithubAccessTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

type GithubUserResponse struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

type GithubOauth struct{}

func (o *GithubOauth) LoginURL(cfg *OAuthConfig) string {
	return fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s", cfg.ClientID)
}

func (o *GithubOauth) RequestAccessToken(code string, cfg *OAuthConfig) (string, error) {
	accessTokenReq := GithubAccessTokenRequest{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.Secret,
		Code:         code,
	}

	data, err := json.Marshal(accessTokenReq)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tokenData := AccessTokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&tokenData)
	if err != nil {
		return "", err
	}

	return tokenData.AccessToken, nil
}

func (o *GithubOauth) RequestUserInfo(access_token string, cfg *OAuthConfig) (*OauthUserInfo, error) {
	// fetch user data
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+access_token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	userData := GithubUserResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		return nil, err
	}

	external_id := strconv.Itoa(userData.ID)
	info := OauthUserInfo{
		Provider:   ProviderTypeGithub,
		Name:       userData.Login,
		ExternalID: external_id,
		AvatarURL:  fmt.Sprintf("https://github.com/%s.png", userData.Login),
	}

	return &info, nil
}
