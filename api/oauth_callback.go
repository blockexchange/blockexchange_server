package api

import (
	"blockexchange/core"
	"blockexchange/types"
	"net/http"
	"time"
)

func (api *Api) OauthCallback(w http.ResponseWriter, user *types.User, new_user bool) error {
	perms := core.GetPermissions(user, true)
	dur := time.Duration(7 * 24 * time.Hour)

	token, err := api.core.CreateJWT(user, perms, dur)
	if err != nil {
		return err
	}
	api.core.SetClaims(w, token, dur)

	return nil
}
