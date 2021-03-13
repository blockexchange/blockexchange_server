package web

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/types"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type TokenApi struct {
	AccessTokenRepo db.AccessTokenRepository
	UserRepo        db.UserRepository
}

func (api TokenApi) PostLogin(w http.ResponseWriter, r *http.Request) {
	login := types.Login{}
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		SendError(w, "decode: "+err.Error())
		return
	}

	logrus.WithFields(logrus.Fields{
		"Username": login.Username,
		"Password": login.Password,
		"Token":    login.Token,
	}).Debug("POST /api/token")

	user, err := api.UserRepo.GetUserByName(login.Username)
	if err != nil {
		SendError(w, "user: "+err.Error())
		return
	}
	if user == nil {
		SendError(w, "User not found")
		return
	}

	if login.Password != "" {
		// login with username / password
		err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(login.Password))
		if err != nil {
			SendError(w, err.Error())
			return
		}

		permissions := []types.JWTPermission{
			types.JWTPermissionUpload,
			types.JWTPermissionOverwrite,
			types.JWTPermissionManagement,
		}

		exp := time.Now().Unix() + (3600 * 24 * 180)
		token, err := core.CreateJWT(user, permissions, exp)
		if err != nil {
			SendError(w, err.Error())
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))

	} else if login.Token != "" {
		// login with token
		access_token, err := api.AccessTokenRepo.GetAccessTokenByTokenAndUserID(login.Token, user.ID)
		if err != nil {
			SendError(w, err.Error())
			return
		}

		if access_token == nil {
			SendError(w, "access token not found")
			return
		}

		if access_token.Expires < (time.Now().Unix() * 1000) {
			SendError(w, "token expired")
			return
		}

		permissions := []types.JWTPermission{
			types.JWTPermissionUpload,
			types.JWTPermissionOverwrite,
		}
		token, err := core.CreateJWT(user, permissions, int64(access_token.Expires/1000))
		if err != nil {
			SendError(w, err.Error())
			return
		}
		api.AccessTokenRepo.IncrementAccessTokenUseCount(access_token.ID)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))

	} else {
		SendError(w, "Empty password/access_token not allowed")
	}

}
