package web

import (
	"blockexchange/core"
	"blockexchange/types"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (api Api) GetAccessTokens(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionManagement) {
		return
	}
	tokens, err := api.AccessTokenRepo.GetAccessTokensByUserID(ctx.Token.UserID)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokens)
}

func (api Api) PostAccessToken(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionManagement) {
		return
	}
	accessToken := types.AccessToken{}
	err := json.NewDecoder(r.Body).Decode(&accessToken)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	accessToken.UserID = ctx.Token.UserID
	accessToken.Token = core.CreateToken(6)
	accessToken.Created = time.Now().Unix() * 1000

	err = api.AccessTokenRepo.CreateAccessToken(&accessToken)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

func (api Api) DeleteAccessToken(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionManagement) {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendError(w, 405, err.Error())
		return
	}

	api.AccessTokenRepo.RemoveAccessToken(int64(id), ctx.Token.UserID)
	w.WriteHeader(http.StatusOK)
}
