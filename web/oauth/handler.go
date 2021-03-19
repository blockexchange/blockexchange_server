package oauth

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/types"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type OauthUserInfo struct {
	Name       string
	Email      string
	ExternalID string
	Type       string
}

type OauthImplementation interface {
	RequestAccessToken(code string) (string, error)
	RequestUserInfo(access_token string) (*OauthUserInfo, error)
}

type OauthHandler struct {
	Impl            OauthImplementation
	AccessTokenRepo db.AccessTokenRepository
	UserRepo        db.UserRepository
}

func NewHandler(impl OauthImplementation, ur db.UserRepository, atr db.AccessTokenRepository) *OauthHandler {
	return &OauthHandler{
		Impl:            impl,
		AccessTokenRepo: atr,
		UserRepo:        ur,
	}
}

func SendError(w http.ResponseWriter, code int, message string) {
	logrus.Trace("web.SendError: " + message)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(types.ErrorResponse{Message: message})
}

func SendJson(w http.ResponseWriter, o interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(o)
}

func (h *OauthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	list := r.URL.Query()["code"]
	if len(list) == 0 {
		SendError(w, 404, "no code found")
		return
	}

	code := list[0]

	access_token, err := h.Impl.RequestAccessToken(code)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	info, err := h.Impl.RequestUserInfo(access_token)
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
	user, err := h.UserRepo.GetUserByName(info.Name)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if user != nil && user.Type != info.Type {
		// assign pseudo-random alternative name
		info.Name = info.Name + "_" + core.CreateToken(6)
	}

	// fetch by external id
	user, err = h.UserRepo.GetUserByExternalId(info.ExternalID)
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
			Hash:       "",
			Mail:       &info.Email,
			ExternalID: &info.ExternalID,
		}
		err = h.UserRepo.CreateUser(user)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		err = h.AccessTokenRepo.CreateAccessToken(&types.AccessToken{
			Name:    "default",
			Created: time.Now().Unix() * 1000,
			Expires: (time.Now().Unix() + (3600 * 24 * 7 * 4)) * 1000,
			Token:   core.CreateToken(6),
			UserID:  user.ID,
		})
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
	}

	permissions := []types.JWTPermission{
		types.JWTPermissionUpload,
		types.JWTPermissionOverwrite,
		types.JWTPermissionManagement,
	}
	exp := time.Now().Unix() + (3600 * 24 * 180)
	token, err := core.CreateJWT(user, permissions, exp)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	target := os.Getenv("BASE_URL") + "/#/oauth/" + token
	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
}
