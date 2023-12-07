package oauth

import (
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"errors"
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
	Email string `json:"email"`
}

type GithubUserMail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

type GithubOauth struct{}

func (o *GithubOauth) RequestAccessToken(code, baseurl string, cfg *types.OAuthConfig) (string, error) {
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

func (o *GithubOauth) RequestUserInfo(access_token string, cfg *types.OAuthConfig) (*OauthUserInfo, error) {
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

	// fetch mails
	req, err = http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+access_token)

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code in email-response: %d", resp.StatusCode)
	}

	mails := []GithubUserMail{}
	err = json.NewDecoder(resp.Body).Decode(&mails)
	if err != nil {
		return nil, err
	}

	// fetch primary mail
	primary_mail := ""
	for _, mail := range mails {
		if mail.Primary && mail.Verified {
			primary_mail = mail.Email
		}
	}

	if primary_mail == "" {
		return nil, errors.New("no primary and verified email address found")
	}

	external_id := strconv.Itoa(userData.ID)
	info := OauthUserInfo{
		Name:       userData.Login,
		Email:      primary_mail,
		ExternalID: external_id,
	}

	return &info, nil
}
