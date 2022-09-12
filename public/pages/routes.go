package pages

import (
	"blockexchange/controller"
	"blockexchange/core"
	"blockexchange/public/oauth"
	"blockexchange/public/pages/schema"
	"blockexchange/types"

	"github.com/gorilla/mux"
)

func SetupRoutes(ctrl *controller.Controller, r *mux.Router, cfg *core.Config) {
	r.HandleFunc("/", ctrl.Handler("./", Index))
	r.HandleFunc("/login", ctrl.Handler("./", Login))
	r.HandleFunc("/schema/{username}", ctrl.Handler("../../", schema.UserSchema))
	r.HandleFunc("/schema/{username}/{schemaname}", ctrl.Handler("../../../", schema.Schema))
	r.HandleFunc("/schema/{username}/{schemaname}/edit", ctrl.Handler("../../../", schema.SchemaEdit, types.JWTPermissionManagement))
	r.HandleFunc("/schema/{username}/{schemaname}/edit-tags", ctrl.Handler("../../../", schema.SchemaTagEdit, types.JWTPermissionManagement))
	r.HandleFunc("/users", ctrl.Handler("./", Users))
	r.HandleFunc("/search", ctrl.Handler("./", Search))
	r.HandleFunc("/mod", ctrl.Handler("./", Mod))
	r.HandleFunc("/profile", ctrl.Handler("./", Profile))
	r.HandleFunc("/tags", ctrl.Handler("./", Tags, types.JWTPermissionAdmin))
	r.HandleFunc("/about", ctrl.Handler("./", About))

	if cfg.DiscordOAuthConfig != nil {
		r.Handle("/api/oauth_callback/discord", oauth.NewHandler(&oauth.DiscordOauth{}, cfg, ctrl.UserRepo, ctrl.AccessTokenRepo, ctrl))
	}
	if cfg.GithubOAuthConfig != nil {
		r.Handle("/api/oauth_callback/github", oauth.NewHandler(&oauth.GithubOauth{}, cfg, ctrl.UserRepo, ctrl.AccessTokenRepo, ctrl))
	}
	if cfg.MesehubOAuthConfig != nil {
		r.Handle("/api/oauth_callback/mesehub", oauth.NewHandler(&oauth.MesehubOauth{}, cfg, ctrl.UserRepo, ctrl.AccessTokenRepo, ctrl))
	}
}
