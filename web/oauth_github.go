package web

import (
	"blockexchange/core"
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

func (api Api) OauthGithub(w http.ResponseWriter, r *http.Request) {
	list := r.URL.Query()["code"]
	if len(list) == 0 {
		SendError(w, 404, "no code found")
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
		SendError(w, 500, err.Error())
		return
	}

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(data))
	if err != nil {
		SendError(w, 500, err.Error())
		return

	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	tokenData := types.GithubAccessTokenRespone{}
	err = json.NewDecoder(resp.Body).Decode(&tokenData)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	fmt.Println(fmt.Sprintf("AccessCode: %s", tokenData.AccessToken))

	req, err = http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)

	resp, err = client.Do(req)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	userData := types.GithubUserResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	fmt.Println(userData)
	user, err := api.UserRepo.GetUserByExternalId(strconv.Itoa(userData.ID))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if user == nil {
		logrus.Debug("creating new user")
		external_id := strconv.Itoa(userData.ID)
		user = &types.User{
			Created:    time.Now().Unix() * 1000,
			Name:       userData.Login,
			Type:       types.UserTypeGithub,
			Hash:       "",
			Mail:       &userData.Email,
			ExternalID: &external_id,
		}
		err = api.UserRepo.CreateUser(user)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		err = api.AccessTokenRepo.CreateAccessToken(&types.AccessToken{
			Name:    "default",
			Created: time.Now().Unix() * 1000,
			Expires: (time.Now().Unix() + (3600 * 24 * 7 * 4)) * 1000,
			Token:   core.CreateToken(6),
			UserID:  user.ID,
		})
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
	}
	permissions := []types.JWTPermission{
		types.JWTPermissionUpload,
		types.JWTPermissionOverwrite,
		types.JWTPermissionManagement,
	}
	exp := time.Now().Unix() + (3600 * 24 * 180)
	token, err := core.CreateJWT(user, permissions, exp)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	target := os.Getenv("BASE_URL") + "/#/oauth/" + token
	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
}
