package web

import (
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func sendError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(types.ErrorResponse{Message: message})
}

func OauthGithub(w http.ResponseWriter, r *http.Request) {
	list := r.URL.Query()["code"]
	if len(list) == 0 {
		sendError(w, "no code found")
		return
	}

	code := list[0]
	fmt.Println(fmt.Sprintf("Code: %s", code))

	accessTokenReq := types.GithubAccessTokenRequest{
		ClientID:     os.Getenv("GITHUB_APP_ID"),
		ClientSecret: os.Getenv("GITHUB_APP_SECRET"),
		Code:         code,
	}

	data, err := json.Marshal(accessTokenReq)
	if err != nil {
		sendError(w, err.Error())
		return
	}

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(data))
	if err != nil {
		sendError(w, err.Error())
		return

	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		sendError(w, err.Error())
		return
	}

	tokenData := types.GithubAccessTokenRespone{}
	err = json.NewDecoder(resp.Body).Decode(&tokenData)
	if err != nil {
		sendError(w, err.Error())
		return
	}

	fmt.Println(fmt.Sprintf("AccessCode: %s", tokenData.AccessToken))

	req, err = http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		sendError(w, err.Error())
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)

	resp, err = client.Do(req)
	if err != nil {
		sendError(w, err.Error())
		return
	}

	userData := types.GithubUserResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		sendError(w, err.Error())
		return
	}

	fmt.Println(userData)
	//TODO

	w.WriteHeader(http.StatusOK)
}
