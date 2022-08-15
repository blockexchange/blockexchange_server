package pages

import (
	"blockexchange/core"
	"blockexchange/public/oauth"
	"blockexchange/types"

	"github.com/gorilla/mux"
)

func (ctrl *Controller) SetupRoutes(r *mux.Router, cfg *core.Config) {
	r.HandleFunc("/", ctrl.Index)
	r.HandleFunc("/login", ctrl.Login)
	r.HandleFunc("/schema/{username}/{schemaname}", ctrl.Schema)
	r.HandleFunc("/schema/{username}/{schemaname}/edit", ctrl.SchemaEdit)
	r.HandleFunc("/users", ctrl.Users)
	r.HandleFunc("/search", ctrl.Search)
	r.HandleFunc("/mod", ctrl.Mod)
	r.HandleFunc("/profile", ctrl.Secure("./", ctrl.Profile))
	r.HandleFunc("/tags", ctrl.Secure("./", ctrl.Tags, types.JWTPermissionAdmin))
	r.HandleFunc("/about", ctrl.About)

	if cfg.DiscordOAuthConfig != nil {
		r.Handle("/api/oauth_callback/discord", oauth.NewHandler(&oauth.DiscordOauth{}, cfg, ctrl.UserRepo, ctrl.AccessTokenRepo, ctrl.te))
	}
	if cfg.GithubOAuthConfig != nil {
		r.Handle("/api/oauth_callback/github", oauth.NewHandler(&oauth.GithubOauth{}, cfg, ctrl.UserRepo, ctrl.AccessTokenRepo, ctrl.te))
	}
	if cfg.MesehubOAuthConfig != nil {
		r.Handle("/api/oauth_callback/mesehub", oauth.NewHandler(&oauth.MesehubOauth{}, cfg, ctrl.UserRepo, ctrl.AccessTokenRepo, ctrl.te))
	}
}
