package web

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/tmpl"
	"blockexchange/types"
	"blockexchange/web/oauth"
	"blockexchange/web/schema"
	"embed"
	"fmt"
	"time"

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

func prettysize(num int) string {
	if num > (1000 * 1000) {
		return fmt.Sprintf("%d MB", num/(1000*1000))
	} else if num > 1000 {
		return fmt.Sprintf("%d kB", num/(1000))
	} else {
		return fmt.Sprintf("%d bytes", num)
	}
}

func formattime(ts int64) string {
	t := time.UnixMilli(ts)
	return t.Format(time.UnixDate)
}

func (ctx *Context) Setup(r *mux.Router) {
	r.HandleFunc("/", ctx.Index)
	r.HandleFunc("/signup", ctx.Signup)
	r.HandleFunc("/login", ctx.tu.OptionalSecure(ctx.Login))
	r.HandleFunc("/search", ctx.tu.OptionalSecure(ctx.Search))
	r.HandleFunc("/profile", ctx.tu.Secure(ctx.Profile, tmpl.PermissionCheck(types.JWTPermissionManagement)))
	r.HandleFunc("/import", ctx.tu.Secure(ctx.SchemaImport, tmpl.PermissionCheck(types.JWTPermissionManagement)))
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
