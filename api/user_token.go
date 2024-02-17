package api

import (
	"blockexchange/core"
	"blockexchange/types"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (api *Api) GetAccessTokens(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	list, err := api.AccessTokenRepo.GetAccessTokensByUserUID(ctx.Claims.UserUID)
	Send(w, list, err)
}

func (api *Api) CreateAccessToken(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	at := &types.AccessToken{}
	err := json.NewDecoder(r.Body).Decode(at)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	at.Created = time.Now().UnixMilli()
	at.Token = core.CreateToken(6)
	at.UserUID = ctx.Claims.UserUID
	err = api.AccessTokenRepo.CreateAccessToken(at)
	Send(w, at, err)
}

func (api *Api) DeleteAccessToken(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	err := api.AccessTokenRepo.RemoveAccessToken(vars["id"], ctx.Claims.UserUID)
	Send(w, true, err)
}
