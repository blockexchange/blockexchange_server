package oauth

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/types"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type OauthHandler struct {
	Impl     OauthImplementation
	UserRepo *db.UserRepository
	Core     *core.Core
	Config   *types.OAuthConfig
	BaseURL  string
	Type     types.UserType
	Callback SuccessCallback
}

func SendJson(w http.ResponseWriter, o any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(o)
}

func SendError(w http.ResponseWriter, code int, message string) {
	logrus.WithFields(logrus.Fields{
		"code":    code,
		"message": message,
	}).Error("http error")
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(code)
	w.Write([]byte(message))
}

func (h *OauthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list := r.URL.Query()["code"]
	if len(list) == 0 {
		SendError(w, 500, "no code found")
		return
	}

	code := list[0]

	access_token, err := h.Impl.RequestAccessToken(code, h.BaseURL, h.Config)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	info, err := h.Impl.RequestUserInfo(access_token, h.Config)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if info.Email == "" {
		SendError(w, 500, "empty email")
		return
	}

	if info.ExternalID == "" {
		SendError(w, 500, "empty external_id")
		return
	}

	// check if there is already a user by that name
	user, err := h.UserRepo.GetUserByName(info.Name)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if user == nil {
		// check register
		rr := &types.RegisterRequest{
			Name: info.Name,
			Mail: info.Email,
		}
		user, err := h.Core.Register(rr, h.Type)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		logrus.WithFields(logrus.Fields{
			"name":        user.Name,
			"type":        user.Type,
			"mail":        user.Mail,
			"external_id": user.ExternalID,
		}).Debug("created new user")

		err = h.Callback(w, user, true)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

	} else {
		// user already registered, check if they are logging in from the same provider
		if user.Type != h.Type {
			SendError(w, 500, "user already registered from a different provider")
			return
		}

		err = h.Callback(w, user, false)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
	}

	target := h.BaseURL + "/profile"
	http.Redirect(w, r, target, http.StatusSeeOther)
}
