package templateengine

import (
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"sync"
)

type TemplateEngineOptions struct {
	Templates    fs.FS
	TemplateDir  string
	EnableCache  bool
	JWTKey       string
	CookieName   string
	CookiePath   string
	CookieDomain string
	CookieSecure bool
}

type TemplateEngine struct {
	options *TemplateEngineOptions
	cache   map[string]*template.Template
	lock    *sync.RWMutex
}

func NewTemplateEngine(options *TemplateEngineOptions) *TemplateEngine {
	return &TemplateEngine{
		options: options,
		cache:   make(map[string]*template.Template),
		lock:    &sync.RWMutex{},
	}
}

func (te *TemplateEngine) GetTemplate(file string) *template.Template {
	te.lock.RLock()
	t := te.cache[file]
	te.lock.RUnlock()

	if t != nil {
		return t
	}

	// create new template
	var f fs.FS = te.options.Templates
	if !te.options.EnableCache {
		// use live fs
		f = os.DirFS(te.options.TemplateDir)
	}

	tmpl := template.Must(template.New("").Option("missingkey=error").ParseFS(f, "layout.html", "components/*.html"))
	t = template.Must(tmpl.ParseFS(f, file))

	if te.options.EnableCache {
		// cache template
		te.lock.Lock()
		te.cache[file] = t
		te.lock.Unlock()
	}

	return t
}

func (te *TemplateEngine) Execute(file string, w http.ResponseWriter, r *http.Request, baseUrl string, data any) error {
	t := te.GetTemplate(file)
	claims, err := te.GetClaims(r)
	if err != nil {
		return te.ExecuteError(w, r, baseUrl, 500, err)
	}

	rd := &RenderData{
		BaseURL: baseUrl,
		Claims:  claims,
		Data:    data,
	}

	return t.ExecuteTemplate(w, "layout", rd)
}

func (te *TemplateEngine) ExecuteError(w http.ResponseWriter, r *http.Request, baseUrl string, statuscode int, err error) error {
	w.WriteHeader(statuscode)
	t := te.GetTemplate("pages/error.html")
	rd := &RenderData{
		BaseURL: baseUrl,
		Data:    err,
	}
	return t.ExecuteTemplate(w, "layout", rd)
}
