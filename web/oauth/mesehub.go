package oauth

import (
	"blockexchange/core"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type MesehubUserResponse struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}

type MesehubOauth struct {
}

func (o *MesehubOauth) RequestAccessToken(code string, cfg *core.Config) (string, error) {
	accessTokenReq := make(map[string]string)
	accessTokenReq["client_id"] = cfg.MesehubOAuthConfig.ClientID
	accessTokenReq["client_secret"] = cfg.MesehubOAuthConfig.Secret
	accessTokenReq["code"] = code
	accessTokenReq["grant_type"] = "authorization_code"
	accessTokenReq["redirect_uri"] = cfg.BaseURL + "/api/oauth_callback/mesehub"

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
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected response-status: %d", resp.StatusCode)
	}

	tokenData := AccessTokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&tokenData)
	if err != nil {
		return "", err
	}

	return tokenData.AccessToken, nil
}

func (o *MesehubOauth) RequestUserInfo(access_token string, cfg *core.Config) (*OauthUserInfo, error) {
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
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected user-info response-status: %d", resp.StatusCode)
	}

	userData := MesehubUserResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		return nil, err
	}

	fmt.Println(userData)
	external_id := strconv.Itoa(userData.ID)
	info := OauthUserInfo{
		Name:       userData.Login,
		Type:       types.UserTypeMesehub,
		Email:      userData.Email,
		ExternalID: external_id,
	}

	return &info, nil
}
