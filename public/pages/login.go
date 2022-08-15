package pages

import (
	"blockexchange/core"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginModel struct {
	Username         string
	Password         string
	LoginError       bool
	Config           *core.Config
	DiscordOauthLink string
	GithubOauthLink  string
	MesehubOauthLink string
}

func (ctrl *Controller) Login(w http.ResponseWriter, r *http.Request) {
	lm := &LoginModel{
		Config: ctrl.cfg,
	}

	if ctrl.cfg.GithubOAuthConfig != nil {
		lm.GithubOauthLink = fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s", ctrl.cfg.GithubOAuthConfig.ClientID)
	}
	if ctrl.cfg.DiscordOAuthConfig != nil {
		lm.DiscordOauthLink = fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify%%20email",
			ctrl.cfg.DiscordOAuthConfig.ClientID,
			url.QueryEscape(ctrl.cfg.BaseURL+"/api/oauth_callback/discord"),
		)
	}
	if ctrl.cfg.MesehubOAuthConfig != nil {
		lm.MesehubOauthLink = fmt.Sprintf("https://git.minetest.land/login/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=STATE",
			ctrl.cfg.MesehubOAuthConfig.ClientID,
			url.QueryEscape(ctrl.cfg.BaseURL+"/api/oauth_callback/mesehub"),
		)
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			ctrl.te.ExecuteError(w, r, "./", 500, err)
			return
		}

		switch r.Form.Get("action") {
		case "login":
			lm.Username = r.Form.Get("username")
			lm.Password = r.Form.Get("password")

			user, err := ctrl.UserRepo.GetUserByName(lm.Username)
			if err != nil {
				ctrl.te.ExecuteError(w, r, "./", 500, err)
				return
			}

			if user == nil {
				lm.LoginError = true
				w.WriteHeader(401)
				break
			}

			err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(lm.Password))
			if err != nil {
				lm.LoginError = true
				w.WriteHeader(401)
				break
			}

			permissions := core.GetPermissions(user, true)
			dur := time.Duration(7 * 24 * time.Hour)
			token, err := core.CreateJWT(user, permissions, dur)
			if err != nil {
				ctrl.te.ExecuteError(w, r, "./", 500, err)
				return
			}

			ctrl.te.SetToken(w, token, dur)
			http.Redirect(w, r, "login", http.StatusSeeOther)

		case "logout":
			ctrl.te.RemoveToken(w)
			http.Redirect(w, r, "login", http.StatusSeeOther)
			return
		}
	}

	ctrl.te.Execute("pages/login.html", w, r, "./", lm)
}
