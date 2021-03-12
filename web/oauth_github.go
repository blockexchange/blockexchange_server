package web

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
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
	user, err := db.GetUserByExternalId(strconv.Itoa(userData.ID))
	if err != nil {
		sendError(w, err.Error())
		return
	}

	if user == nil {
		logrus.Debug("creating new user")
		user = &types.User{
			Created:    time.Now().Unix() * 1000,
			Name:       userData.Login,
			Type:       types.UserTypeGithub,
			Hash:       "",
			Mail:       userData.Email,
			ExternalID: strconv.Itoa(userData.ID),
		}
		err = db.CreateUser(user)
		if err != nil {
			sendError(w, err.Error())
			return
		}

		err = db.CreateAccessToken(&types.AccessToken{
			Name:    "default",
			Created: time.Now().Unix() * 1000,
			Expires: (time.Now().Unix() + (3600 * 24 * 7 * 4)) * 1000,
			Token:   core.CreateToken(6),
			UserID:  user.ID,
		})
		if err != nil {
			sendError(w, err.Error())
			return
		}
	}
	fmt.Println(user)
	permissions := []types.JWTPermission{
		types.JWTPermissionUpload,
		types.JWTPermissionOverwrite,
		types.JWTPermissionManagement,
	}
	token, err := core.CreateJWT(user, permissions)
	if err != nil {
		sendError(w, err.Error())
		return
	}

	target := os.Getenv("BASE_URL") + "/#/oauth/" + token
	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
}
