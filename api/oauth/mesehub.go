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

type MesehubUserResponse struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}

type MesehubOauth struct{}

func (o *MesehubOauth) RequestAccessToken(code, baseurl string, cfg *types.OAuthConfig) (string, error) {
	accessTokenReq := make(map[string]string)
	accessTokenReq["client_id"] = cfg.ClientID
	accessTokenReq["client_secret"] = cfg.Secret
	accessTokenReq["code"] = code
	accessTokenReq["grant_type"] = "authorization_code"
	accessTokenReq["redirect_uri"] = baseurl + "/oauth_callback/mesehub"

	data, err := json.Marshal(accessTokenReq)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://git.minetest.land/login/oauth/access_token", bytes.NewBuffer(data))
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

func (o *MesehubOauth) RequestUserInfo(access_token string, cfg *types.OAuthConfig) (*OauthUserInfo, error) {
	req, err := http.NewRequest("GET", "https://git.minetest.land/api/v1/user", nil)
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
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code in response: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	userData := MesehubUserResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		return nil, err
	}

	// fetch mails
	req, err = http.NewRequest("GET", "https://git.minetest.land/api/v1/user/emails", nil)
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
