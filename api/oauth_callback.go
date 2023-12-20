package api

import (
	"blockexchange/core"
	"blockexchange/oauth"
	"blockexchange/types"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func (api *Api) OauthCallback(w http.ResponseWriter, r *http.Request, user_info *oauth.OauthUserInfo) error {

	user, err := api.UserRepo.GetUserByExternalIdAndType(user_info.ExternalID, types.UserType(user_info.Provider))
	if err != nil {
		SendError(w, 500, err.Error())
		return nil
	}

	if user == nil {
		user, err = api.core.RegisterOauth(user_info.Name, user_info.ExternalID, types.UserType(user_info.Provider))
		if err != nil {
			SendError(w, 500, err.Error())
			return nil
		}

		logrus.WithFields(logrus.Fields{
			"name":        user.Name,
			"type":        user.Type,
			"external_id": user.ExternalID,
		}).Debug("created new user")
	}

	perms := core.GetPermissions(user, true)
	dur := time.Duration(7 * 24 * time.Hour)

	token, err := api.core.CreateJWT(user, perms, dur)
	if err != nil {
		return err
	}
	api.core.SetClaims(w, token, dur)

	target := api.cfg.BaseURL + "/profile"
	http.Redirect(w, r, target, http.StatusSeeOther)

	return nil
}
