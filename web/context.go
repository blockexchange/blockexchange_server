package web

import (
	"blockexchange/types"
	"blockexchange/web/oauth"
	"embed"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

//go:embed *
var Files embed.FS

func (ctx *Context) CreateTemplate(pagename string, r *http.Request) *template.Template {
	funcs := template.FuncMap{
		"BaseURL": func() string { return ctx.BaseURL },
		"Claims": func() *types.Claims {
			c, _ := ctx.GetClaims(r)
			return c
		},
		"prettysize": prettysize,
		"formattime": formattime,
	}
	return template.Must(template.New("").Funcs(funcs).ParseFS(Files, "components/*.html", pagename))
}

func (ctx *Context) StaticPage(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := ctx.CreateTemplate(name, r)
		t.ExecuteTemplate(w, "layout", nil)
	}
}

func (ctx *Context) Setup(r *mux.Router) {
	r.HandleFunc("/", ctx.Index)
	r.HandleFunc("/login", ctx.OptionalSecure(ctx.Login))
	r.HandleFunc("/search", ctx.OptionalSecure(ctx.Search))
	r.HandleFunc("/profile", ctx.Secure(ctx.Profile, permissionCheck(types.JWTPermissionManagement)))
	r.HandleFunc("/users", ctx.Users)
	r.HandleFunc("/mod", ctx.StaticPage("mod.html"))
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
