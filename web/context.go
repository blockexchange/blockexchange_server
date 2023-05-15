package web

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

//go:embed *
var Files embed.FS

func createTemplate(pagename string) *template.Template {
	funcs := template.FuncMap{
		"getPrefixURL": func() string { return "/" },
	}
	return template.Must(template.New("").Funcs(funcs).ParseFS(Files, "components/*.html", pagename))
}

func (ctx *Context) StaticPage(name string) http.HandlerFunc {
	t := createTemplate(name)
	return ctx.OptionalSecure(func(w http.ResponseWriter, r *http.Request, c *Claims) {
		t.ExecuteTemplate(w, "layout", map[string]any{
			"Claims": c,
		})
	})
}

func (ctx *Context) Setup() {
	r := mux.NewRouter()

	r.HandleFunc("/", ctx.StaticPage("index.html"))
	r.PathPrefix("/assets").Handler(http.FileServer(http.FS(Files)))

	r.NotFoundHandler = ctx.NotFound()
	http.Handle("/", r)
}
