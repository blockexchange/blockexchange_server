package api

import (
	"blockexchange/core"
	"blockexchange/types"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func (api *Api) RequestToken(w http.ResponseWriter, r *http.Request) {
	login := types.Login{}
	err := json.NewDecoder(r.Body).Decode(&login)
	if SendIfError(w, err) {
		return
	}

	logrus.WithFields(logrus.Fields{
		"Username": login.Username,
	}).Debug("POST /api/token")

	user, err := api.UserRepo.GetUserByName(login.Username)
	if SendIfError(w, err) {
		return
	}
	if user == nil {
		SendError(w, 404, "User not found")
		return
	}

	if login.Token != "" {
		// login with token
		access_token, err := api.AccessTokenRepo.GetAccessTokenByTokenAndUserUID(login.Token, user.UID)
		if SendIfError(w, err) {
			return
		}

		if access_token == nil {
			SendError(w, 404, "access token not found")
			return
		}

		if access_token.Expires < time.Now().UnixMilli() {
			SendError(w, 401, "token expired")
			return
		}

		permissions := core.GetPermissions(user, false)
		token, err := api.core.CreateJWT(user, permissions, time.Until(time.UnixMilli(access_token.Expires)))
		if SendIfError(w, err) {
			return
		}
		api.AccessTokenRepo.IncrementAccessTokenUseCount(access_token.UID)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))

	} else {
		SendError(w, 405, "Empty access_token not allowed")
	}

}
