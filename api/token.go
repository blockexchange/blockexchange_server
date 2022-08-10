package api

import (
	"blockexchange/core"
	"blockexchange/types"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (api Api) PostLogin(w http.ResponseWriter, r *http.Request) {
	login := types.Login{}
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	logrus.WithFields(logrus.Fields{
		"Username": login.Username,
		"Password": login.Password,
		"Token":    login.Token,
	}).Debug("POST /api/token")

	user, err := api.UserRepo.GetUserByName(login.Username)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if user == nil {
		SendError(w, 404, "User not found")
		return
	}

	if login.Password != "" {
		// login with username / password
		err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(login.Password))
		if err != nil {
			SendError(w, 401, err.Error())
			return
		}

		permissions := core.GetPermissions(user, true)
		if user.Role == types.UserRoleAdmin {
			// admin user
			permissions = append(permissions, types.JWTPermissionAdmin)
		}

		exp := time.Now().Unix() + (3600 * 24 * 180)
		token, err := core.CreateJWT(user, permissions, exp)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))

	} else if login.Token != "" {
		// login with token
		access_token, err := api.AccessTokenRepo.GetAccessTokenByTokenAndUserID(login.Token, user.ID)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		if access_token == nil {
			SendError(w, 404, "access token not found")
			return
		}

		if access_token.Expires < (time.Now().Unix() * 1000) {
			SendError(w, 401, "token expired")
			return
		}

		permissions := core.GetPermissions(user, false)
		token, err := core.CreateJWT(user, permissions, int64(access_token.Expires/1000))
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
		api.AccessTokenRepo.IncrementAccessTokenUseCount(access_token.ID)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))

	} else {
		SendError(w, 405, "Empty password/access_token not allowed")
	}

}
