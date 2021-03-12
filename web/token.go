package web

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/types"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func PostLogin(w http.ResponseWriter, r *http.Request) {
	login := types.Login{}
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		SendError(w, err.Error())
	}

	user, err := db.GetUserByName(login.Username)
	if err != nil {
		SendError(w, err.Error())
		return
	}

	if login.Password != "" {
		// login with username / password
		bcrypt.CompareHashAndPassword(nil, nil)
		// TODO
	} else if login.Token != "" {
		// login with token
		access_token, err := db.GetAccessTokenByTokenAndUserID(login.Token, user.ID)
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
		db.IncrementAccessTokenUseCount(access_token.ID)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(token))

	} else {
		SendError(w, "Empty password/access_token not allowed")
	}

}
