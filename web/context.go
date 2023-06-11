package web

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/tmpl"
	"blockexchange/types"
	"blockexchange/web/oauth"
	"blockexchange/web/schema"
	"embed"

	"github.com/gorilla/mux"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

//go:embed *
var Files embed.FS

type Context struct {
	Repos  *db.Repositories
	Config *core.Config
	tu     *tmpl.TemplateUtil
}

func (ctx *Context) Setup(r *mux.Router) {
	r.HandleFunc("/", ctx.Index)
	r.HandleFunc("/signup", ctx.Signup)
	r.HandleFunc("/login", ctx.tu.OptionalSecure(ctx.Login))
	r.HandleFunc("/search", ctx.tu.OptionalSecure(ctx.Search))
	r.HandleFunc("/profile", ctx.tu.Secure(ctx.Profile, tmpl.PermissionCheck(types.JWTPermissionManagement)))
	r.HandleFunc("/import", ctx.tu.Secure(ctx.SchemaImport, tmpl.PermissionCheck(types.JWTPermissionManagement)))
	r.HandleFunc("/tags", ctx.tu.Secure(ctx.Tags, tmpl.PermissionCheck(types.JWTPermissionAdmin)))
	r.HandleFunc("/users", ctx.Users)
	r.HandleFunc("/mod", ctx.tu.StaticPage("mod.html"))
	r.PathPrefix("/assets").Handler(statigz.FileServer(Files, brotli.AddEncoding))

	if ctx.Config.DiscordOAuthConfig != nil {
		r.Handle("/api/oauth_callback/discord", NewHandler(&oauth.DiscordOauth{}, ctx))
	}
	if ctx.Config.GithubOAuthConfig != nil {
		r.Handle("/api/oauth_callback/github", NewHandler(&oauth.GithubOauth{}, ctx))
	}
	if ctx.Config.MesehubOAuthConfig != nil {
		r.Handle("/api/oauth_callback/mesehub", NewHandler(&oauth.MesehubOauth{}, ctx))
	}

	sc := schema.NewSchemaContext(ctx.tu, *ctx.Repos, ctx.Config.BaseURL)
	sc.Setup(r)

	r.NotFoundHandler = ctx.NotFound()
}
