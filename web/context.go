package web

import (
	"blockexchange/types"
	"embed"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

//go:embed *
var Files embed.FS

func (ctx *Context) CreateTemplate(pagename string) *template.Template {
	funcs := template.FuncMap{
		"BaseURL":    func() string { return ctx.BaseURL },
		"prettysize": prettysize,
		"formattime": formattime,
	}
	return template.Must(template.New("").Funcs(funcs).ParseFS(Files, "components/*.html", pagename))
}

func (ctx *Context) StaticPage(name string) http.HandlerFunc {
	t := ctx.CreateTemplate(name)
	return ctx.OptionalSecure(func(w http.ResponseWriter, r *http.Request, c *types.Claims) {
		t.ExecuteTemplate(w, "layout", map[string]any{
			"Claims": c,
		})
	})
}

func (ctx *Context) Setup(r *mux.Router) {
	r.HandleFunc("/", ctx.OptionalSecure(ctx.Index))
	r.HandleFunc("/login", ctx.OptionalSecure(ctx.Login))
	r.HandleFunc("/mod", ctx.StaticPage("mod.html"))
	r.PathPrefix("/assets").Handler(statigz.FileServer(Files, brotli.AddEncoding))

	ctx.error_template = ctx.CreateTemplate("error.html")
	r.NotFoundHandler = ctx.NotFound()
}
