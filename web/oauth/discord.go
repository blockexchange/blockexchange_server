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

type DiscordResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type DiscordOauth struct {
}

func (o *DiscordOauth) RequestAccessToken(code string) (string, error) {
	accessTokenReq := make(map[string]interface{})
	accessTokenReq["client_id"] = os.Getenv("DISCORD_APP_ID")
	accessTokenReq["client_secret"] = os.Getenv("DISCORD_APP_SECRET")
	accessTokenReq["redirect_uri"] = os.Getenv("BASE_URL") + "/api/oauth_callback/discord"
	accessTokenReq["code"] = code
	accessTokenReq["grant_type"] = "authorization_code"
	accessTokenReq["scope"] = "identify email connections"

	data, err := json.Marshal(accessTokenReq)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://discord.com/api/oauth2/token", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//TODO: query string shenanigans
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	tokenData := AccessTokenRespone{}
	err = json.NewDecoder(resp.Body).Decode(&tokenData)
	if err != nil {
		return "", err
	}

	fmt.Println(fmt.Sprintf("AccessCode: %s", tokenData.AccessToken))
	return tokenData.AccessToken, nil
}

func (o *DiscordOauth) RequestUserInfo(access_token string) (*OauthUserInfo, error) {
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
	external_id := strconv.Itoa(userData.ID)
	info := OauthUserInfo{
		Name:       userData.Username,
		Type:       types.UserTypeDiscord,
		Email:      userData.Email,
		ExternalID: external_id,
	}

	return &info, nil
}
