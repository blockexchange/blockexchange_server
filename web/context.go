package web

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/types"
	"blockexchange/web/oauth"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

//go:embed *
var Files embed.FS

type Context struct {
	JWTKey       string
	CookieName   string
	CookieDomain string
	CookiePath   string
	CookieSecure bool
	BaseURL      string
	Repos        *db.Repositories
	Config       *core.Config
	tu           *TemplateUtil
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
	ctx.tu = &TemplateUtil{
		Files: Files,
		AddFuncs: func(funcs template.FuncMap, r *http.Request) {
			funcs["BaseURL"] = func() string { return ctx.Config.BaseURL }
			funcs["Claims"] = func() (*types.Claims, error) { return ctx.GetClaims(r) }
			funcs["prettysize"] = prettysize
			funcs["formattime"] = formattime
		},
	}

	r.HandleFunc("/", ctx.Index)
	r.HandleFunc("/signup", ctx.Signup)
	r.HandleFunc("/login", ctx.OptionalSecure(ctx.Login))
	r.HandleFunc("/search", ctx.OptionalSecure(ctx.Search))
	r.HandleFunc("/profile", ctx.Secure(ctx.Profile, permissionCheck(types.JWTPermissionManagement)))
	r.HandleFunc("/import", ctx.Secure(ctx.SchemaImport, permissionCheck(types.JWTPermissionManagement)))
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

	r.NotFoundHandler = ctx.NotFound()
}
