package oauth

import (
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type GithubAccessTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

type GithubAccessTokenRespone struct {
	AccessToken string `json:"access_token"`
}

type GithubUserResponse struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}

type GithubOauth struct {
}

func (o *GithubOauth) RequestAccessToken(code string) (string, error) {
	accessTokenReq := GithubAccessTokenRequest{
		ClientID:     os.Getenv("GITHUB_APP_ID"),
		ClientSecret: os.Getenv("GITHUB_APP_SECRET"),
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

	tokenData := GithubAccessTokenRespone{}
	err = json.NewDecoder(resp.Body).Decode(&tokenData)
	if err != nil {
		return "", err
	}

	fmt.Println(fmt.Sprintf("AccessCode: %s", tokenData.AccessToken))
	return tokenData.AccessToken, nil
}

func (o *GithubOauth) RequestUserInfo(access_token string) (*OauthUserInfo, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
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

	userData := GithubUserResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		return nil, err
	}

	fmt.Println(userData)
	external_id := strconv.Itoa(userData.ID)
	info := OauthUserInfo{
		Name:       userData.Login,
		Type:       types.UserTypeGithub,
		Email:      userData.Email,
		ExternalID: external_id,
	}

	return &info, nil
}
