package api

import (
	"blockexchange/core"
	"blockexchange/types"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (api *Api) DoLogin(w http.ResponseWriter, r *http.Request) {
	login := types.Login{}
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	user, err := api.Repositories.UserRepo.GetUserByName(login.Username)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if user == nil {
		SendError(w, 404, "user not found")
		return
	}
	if user.Type != types.UserTypeLocal {
		SendError(w, 405, "not a local user")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(login.Password))
	if err != nil {
		SendError(w, 401, err.Error())
		return
	}

	permissions := core.GetPermissions(user, true)
	dur := time.Duration(7 * 24 * time.Hour)
	token, err := api.core.CreateJWT(user, permissions, dur)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	api.core.SetClaims(w, token, dur)
	Send(w, core.CreateClaims(user, permissions), nil)
}

func (api *Api) GetLogin(w http.ResponseWriter, r *http.Request) {
	c, err := api.core.GetClaims(r)
	if c == nil {
		Send(w, nil, err)
		return
	}

	if r.URL.Query().Get("renew") == "true" || c.ExpiresAt.Before(time.Now().Add(1*time.Hour)) {
		// renew claims
		user, err := api.UserRepo.GetUserByUID(c.UserUID)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
		dur := time.Duration(7 * 24 * time.Hour)
		permissions := core.GetPermissions(user, true)
		token, err := api.core.CreateJWT(user, permissions, dur)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
		api.core.SetClaims(w, token, dur)
		c = core.CreateClaims(user, permissions)
	}
	Send(w, c, err)
}

func (api *Api) Logout(w http.ResponseWriter, r *http.Request) {
	api.core.RemoveClaims(w)
}
