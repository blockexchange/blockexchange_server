package web

import (
	"blockexchange/core"
	"blockexchange/types"
	"blockexchange/web/oauth"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type OauthHandler struct {
	Impl oauth.OauthImplementation
	ctx  *Context
}

func NewHandler(impl oauth.OauthImplementation, ctx *Context) *OauthHandler {
	return &OauthHandler{
		Impl: impl,
		ctx:  ctx,
	}
}

func SendError(w http.ResponseWriter, code int, message string) {
	logrus.Trace("web.SendError: " + message)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(types.ErrorResponse{Message: message})
}

func SendJson(w http.ResponseWriter, o any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(o)
}

func (h *OauthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list := r.URL.Query()["code"]
	if len(list) == 0 {
		SendError(w, 404, "no code found")
		return
	}

	code := list[0]

	access_token, err := h.Impl.RequestAccessToken(code, h.ctx.Config)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	info, err := h.Impl.RequestUserInfo(access_token, h.ctx.Config)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if info.Name == "" {
		SendError(w, 500, "empty username")
		return
	}

	if info.ExternalID == "" {
		SendError(w, 500, "empty externalid")
		return
	}

	// check if there is already a user by that name
	user, err := h.ctx.Repos.UserRepo.GetUserByName(info.Name)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if user != nil && user.Type != info.Type {
		// assign pseudo-random alternative name
		info.Name = info.Name + "_" + core.CreateToken(6)
	}

	// fetch by external id
	user, err = h.ctx.Repos.UserRepo.GetUserByExternalId(info.ExternalID)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if user == nil {
		logrus.Debug("creating new user")
		user = &types.User{
			Created:    time.Now().Unix() * 1000,
			Name:       info.Name,
			Type:       info.Type,
			Role:       types.UserRoleDefault,
			Hash:       "",
			Mail:       &info.Email,
			ExternalID: &info.ExternalID,
		}
		err = h.ctx.Repos.UserRepo.CreateUser(user)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		err = h.ctx.Repos.AccessTokenRepo.CreateAccessToken(&types.AccessToken{
			Name:    "default",
			Created: time.Now().Unix() * 1000,
			Expires: (time.Now().Unix() + (3600 * 24 * 7 * 4)) * 1000,
			Token:   core.CreateToken(6),
			UserID:  *user.ID,
		})
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
	}

	dur := time.Duration(24 * 180 * time.Hour)
	permissions := core.GetPermissions(user, true)
	token, err := core.CreateJWT(user, permissions, dur)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	h.ctx.SetClaims(w, token, dur)

	target := h.ctx.Config.BaseURL + "/profile"
	http.Redirect(w, r, target, http.StatusSeeOther)
}
