package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
)

type CDBUserResponse struct {
	Username string `json:"username"`
}

type CDBOauth struct{}

func (o *CDBOauth) LoginURL(cfg *OAuthConfig) string {
	return fmt.Sprintf("https://content.minetest.net/oauth/authorize/?response_type=code&client_id=%s&redirect_uri=%s", cfg.ClientID, url.QueryEscape(cfg.CallbackURL))
}

func (o *CDBOauth) RequestAccessToken(code string, cfg *OAuthConfig) (string, error) {
	var data bytes.Buffer
	w := multipart.NewWriter(&data)
	w.WriteField("grant_type", "authorization_code")
	w.WriteField("client_id", cfg.ClientID)
	w.WriteField("client_secret", cfg.Secret)
	w.WriteField("code", code)
	w.Close()

	req, err := http.NewRequest("POST", "https://content.minetest.net/oauth/token/", &data)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status-code: %d", resp.StatusCode)
	}

	tokenData := AccessTokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&tokenData)
	if err != nil {
		return "", err
	}

	return tokenData.AccessToken, nil
}

func (o *CDBOauth) RequestUserInfo(access_token string, cfg *OAuthConfig) (*OauthUserInfo, error) {
	// fetch user data
	req, err := http.NewRequest("GET", "https://content.minetest.net/api/whoami/", nil)
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

	userData := CDBUserResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		return nil, err
	}

	info := OauthUserInfo{
		Provider:   ProviderTypeCDB,
		Name:       userData.Username,
		ExternalID: userData.Username,
	}

	return &info, nil
}
