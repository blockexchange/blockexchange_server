package pages

import (
	"blockexchange/controller"
	"blockexchange/core"
	"blockexchange/public/oauth"
	"blockexchange/types"

	"github.com/gorilla/mux"
)

func SetupRoutes(ctrl *controller.Controller, r *mux.Router, cfg *core.Config) {
	r.HandleFunc("/", ctrl.Handler("./", Index))
	r.HandleFunc("/login", ctrl.Handler("./", Login))
	r.HandleFunc("/schema/{username}/{schemaname}", ctrl.Handler("../../../", Schema))
	r.HandleFunc("/schema/{username}/{schemaname}/edit", ctrl.SecureHandler("../../../", SchemaEdit))
	r.HandleFunc("/users", ctrl.Handler("./", Users))
	r.HandleFunc("/search", ctrl.Handler("./", Search))
	r.HandleFunc("/mod", ctrl.Handler("./", Mod))
	r.HandleFunc("/profile", ctrl.SecureHandler("./", Profile))
	r.HandleFunc("/tags", ctrl.SecureHandler("./", Tags, types.JWTPermissionAdmin))
	r.HandleFunc("/about", ctrl.Handler("./", About))

	if cfg.DiscordOAuthConfig != nil {
		r.Handle("/api/oauth_callback/discord", oauth.NewHandler(&oauth.DiscordOauth{}, cfg, ctrl.UserRepo, ctrl.AccessTokenRepo, ctrl.TemplateEngine()))
	}
	if cfg.GithubOAuthConfig != nil {
		r.Handle("/api/oauth_callback/github", oauth.NewHandler(&oauth.GithubOauth{}, cfg, ctrl.UserRepo, ctrl.AccessTokenRepo, ctrl.TemplateEngine()))
	}
	if cfg.MesehubOAuthConfig != nil {
		r.Handle("/api/oauth_callback/mesehub", oauth.NewHandler(&oauth.MesehubOauth{}, cfg, ctrl.UserRepo, ctrl.AccessTokenRepo, ctrl.TemplateEngine()))
	}
}
