package api

import (
	"blockexchange/core"
	"blockexchange/types"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// longest allowed lifetime of an access token (the ui offers up to one year)
const maxAccessTokenLifetime = 366 * 24 * time.Hour

func (api *Api) GetAccessTokens(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	list, err := api.AccessTokenRepo.GetAccessTokensByUserUID(ctx.Claims.UserUID)
	Send(w, list, err)
}

func (api *Api) CreateAccessToken(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	at := &types.AccessToken{}
	err := json.NewDecoder(r.Body).Decode(at)
	if SendIfError(w, err) {
		return
	}

	now := time.Now()
	if at.Expires <= now.UnixMilli() || at.Expires > now.Add(maxAccessTokenLifetime).UnixMilli() {
		SendError(w, 400, "expires must be a unix timestamp in milliseconds, in the future and at most one year ahead")
		return
	}

	// never trust client-provided identifiers or counters
	at.UID = ""
	at.UseCount = 0
	at.Created = now.UnixMilli()
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
