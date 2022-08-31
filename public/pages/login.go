package pages

import (
	"blockexchange/controller"
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

func Login(rc *controller.RenderContext, r *http.Request) error {
	cfg := rc.Config()
	repos := rc.Repositories()

	lm := &LoginModel{
		Config: cfg,
	}

	if cfg.GithubOAuthConfig != nil {
		lm.GithubOauthLink = fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s", cfg.GithubOAuthConfig.ClientID)
	}
	if cfg.DiscordOAuthConfig != nil {
		lm.DiscordOauthLink = fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify%%20email",
			cfg.DiscordOAuthConfig.ClientID,
			url.QueryEscape(cfg.BaseURL+"/api/oauth_callback/discord"),
		)
	}
	if cfg.MesehubOAuthConfig != nil {
		lm.MesehubOauthLink = fmt.Sprintf("https://git.minetest.land/login/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=STATE",
			cfg.MesehubOAuthConfig.ClientID,
			url.QueryEscape(cfg.BaseURL+"/api/oauth_callback/mesehub"),
		)
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			return err
		}

		switch r.Form.Get("action") {
		case "login":
			lm.Username = r.Form.Get("username")
			lm.Password = r.Form.Get("password")

			user, err := repos.UserRepo.GetUserByName(lm.Username)
			if err != nil {
				return err
			}

			if user == nil {
				lm.LoginError = true
				rc.ResponseWriter().WriteHeader(401)
				break
			}

			err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(lm.Password))
			if err != nil {
				lm.LoginError = true
				rc.ResponseWriter().WriteHeader(401)
				break
			}

			permissions := core.GetPermissions(user, true)
			dur := time.Duration(7 * 24 * time.Hour)
			token, err := core.CreateJWT(user, permissions, dur)
			if err != nil {
				return err
			}

			rc.SetToken(token, dur)
			http.Redirect(rc.ResponseWriter(), r, "login", http.StatusSeeOther)
			return nil

		case "logout":
			rc.RemoveToken()
			http.Redirect(rc.ResponseWriter(), r, "login", http.StatusSeeOther)
			return nil
		}
	}

	return rc.Render("pages/login.html", lm)
}
