package web

import (
	"blockexchange/core"
	"blockexchange/types"
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
	ExternalLogin    bool
	DiscordOauthLink string
	GithubOauthLink  string
	MesehubOauthLink string
	Claims           *types.Claims
}

func (ctx *Context) Login(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	cfg := ctx.Config

	lm := &LoginModel{
		Config: cfg,
		Claims: c,
	}

	if cfg.GithubOAuthConfig != nil {
		lm.ExternalLogin = true
		lm.GithubOauthLink = fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s", cfg.GithubOAuthConfig.ClientID)
	}
	if cfg.DiscordOAuthConfig != nil {
		lm.ExternalLogin = true
		lm.DiscordOauthLink = fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify%%20email",
			cfg.DiscordOAuthConfig.ClientID,
			url.QueryEscape(cfg.BaseURL+"api/oauth_callback/discord"),
		)
	}
	if cfg.MesehubOAuthConfig != nil {
		lm.ExternalLogin = true
		lm.MesehubOauthLink = fmt.Sprintf("https://git.minetest.land/login/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=STATE",
			cfg.MesehubOAuthConfig.ClientID,
			url.QueryEscape(cfg.BaseURL+"api/oauth_callback/mesehub"),
		)
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			ctx.RenderError(w, r, 500, err)
			return
		}

		switch r.Form.Get("action") {
		case "login":
			lm.Username = r.Form.Get("username")
			lm.Password = r.Form.Get("password")

			user, err := ctx.Repos.UserRepo.GetUserByName(lm.Username)
			if err != nil {
				ctx.RenderError(w, r, 500, err)
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
				ctx.RenderError(w, r, 500, err)
				return
			}

			ctx.SetClaims(w, token, dur)
			http.Redirect(w, r, "login", http.StatusSeeOther)
			return

		case "logout":
			ctx.ClearClaims(w)
			http.Redirect(w, r, "login", http.StatusSeeOther)
			return
		}
	}

	t := ctx.CreateTemplate("login.html")
	err := t.ExecuteTemplate(w, "layout", lm)
	if err != nil {
		panic(err)
		//ctx.RenderError(w, r, 500, err)
	}
}
