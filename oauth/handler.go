package oauth

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type OauthHandler struct {
	Impl     OauthImplementation
	Config   *OAuthConfig
	BaseURL  string
	Callback OauthCallback
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

	if info.ExternalID == "" {
		SendError(w, 500, "empty external_id")
		return
	}

	err = h.Callback(w, r, info)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

}
