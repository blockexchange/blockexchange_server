package oauth

import (
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

type CDBAccessTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

type CDBUserResponse struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

type CDBOauth struct{}

func (o *CDBOauth) RequestAccessToken(code, baseurl string, cfg *types.OAuthConfig) (string, error) {
	accessTokenReq := GithubAccessTokenRequest{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.Secret,
		Code:         code,
	}

	data, err := json.Marshal(accessTokenReq)
	if err != nil {
		return "", err
	}

	// TODO
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

func (o *CDBOauth) RequestUserInfo(access_token string, cfg *types.OAuthConfig) (*OauthUserInfo, error) {
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
		Name:       userData.Login,
		ExternalID: external_id,
	}

	return &info, nil
}
