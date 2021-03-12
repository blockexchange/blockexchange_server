package web

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/types"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type AccessTokenApi struct {
	Repo *db.AccessTokenRepository
}

func (api AccessTokenApi) GetAccessTokens(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionManagement) {
		return
	}
	tokens, err := api.Repo.GetAccessTokensByUserID(ctx.Token.UserID)
	if err != nil {
		SendError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokens)
}

func (api AccessTokenApi) PostAccessToken(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionManagement) {
		return
	}
	accessToken := types.AccessToken{}
	err := json.NewDecoder(r.Body).Decode(&accessToken)
	if err != nil {
		SendError(w, err.Error())
		return
	}

	accessToken.UserID = ctx.Token.UserID
	accessToken.Token = core.CreateToken(6)
	accessToken.Created = time.Now().Unix() * 1000

	err = api.Repo.CreateAccessToken(&accessToken)
	if err != nil {
		SendError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

func (api AccessTokenApi) DeleteAccessToken(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionManagement) {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendError(w, err.Error())
		return
	}

	api.Repo.RemoveAccessToken(int64(id), ctx.Token.UserID)
	w.WriteHeader(http.StatusOK)
}
