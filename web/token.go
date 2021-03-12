package web

import (
	"blockexchange/db"
	"blockexchange/types"
	"encoding/json"
	"net/http"
	"time"
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

		//TODO

	} else {
		SendError(w, "Specify a username, password or access_token")
	}

}
