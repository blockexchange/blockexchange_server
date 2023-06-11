package tmpl

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"
)

type TemplateUtil struct {
	Files        embed.FS
	AddFuncs     func(funcs template.FuncMap, r *http.Request)
	JWTKey       string
	CookieName   string
	CookieDomain string
	CookiePath   string
	CookieSecure bool
}

func (tu *TemplateUtil) CreateTemplate(pagename string, r *http.Request) *template.Template {
	funcs := template.FuncMap{
		"Claims": func() (any, error) { return tu.GetClaims(r) },
	}
	tu.AddFuncs(funcs, r)
	return template.Must(template.New("").Funcs(funcs).ParseFS(tu.Files, "components/*.html", pagename))
}

func (tu *TemplateUtil) StaticPage(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := tu.CreateTemplate(name, r)
		t.ExecuteTemplate(w, "layout", nil)
	}
}

func (tu *TemplateUtil) ExecuteTemplate(w http.ResponseWriter, r *http.Request, name string, data any) {
	t := tu.CreateTemplate(name, r)
	buf := bytes.NewBuffer([]byte{})
	err := t.ExecuteTemplate(buf, "layout", data)
	if err != nil {
		tu.RenderError(w, r, 500, err)
	} else {
		w.Write(buf.Bytes())
	}
}

func (tu *TemplateUtil) RenderError(w http.ResponseWriter, r *http.Request, code int, err error) {
	w.WriteHeader(code)
	t := tu.CreateTemplate("error.html", r)
	t.ExecuteTemplate(w, "layout", err)
}
