package templateengine

import (
	"bytes"
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

func (te *TemplateEngine) GetTemplate(file string) (*template.Template, error) {
	te.lock.RLock()
	t := te.cache[file]
	te.lock.RUnlock()

	if t != nil {
		return t, nil
	}

	// create new template
	var f fs.FS = te.options.Templates
	if !te.options.EnableCache {
		// use live fs
		f = os.DirFS(te.options.TemplateDir)
	}

	tmpl, err := template.New("").Option("missingkey=error").ParseFS(f, "layout.html", "components/*.html")
	if err != nil {
		return nil, err
	}

	tmpl, err = tmpl.ParseFS(f, file)
	if err != nil {
		return nil, err
	}

	if te.options.EnableCache {
		// cache template
		te.lock.Lock()
		te.cache[file] = tmpl
		te.lock.Unlock()
	}

	return tmpl, nil
}

func (te *TemplateEngine) Execute(file string, w http.ResponseWriter, r *http.Request, baseUrl string, data any) error {
	t, err := te.GetTemplate(file)
	if err != nil {
		return te.ExecuteError(w, r, baseUrl, 500, err)
	}

	claims, err := te.GetClaims(r)
	if err != nil {
		return te.ExecuteError(w, r, baseUrl, 500, err)
	}

	rd := &RenderData{
		BaseURL: baseUrl,
		Claims:  claims,
		Data:    data,
	}

	buf := bytes.Buffer{}

	err = t.ExecuteTemplate(&buf, "layout", rd)
	if err != nil {
		te.ExecuteError(w, r, baseUrl, 500, err)
		return err
	}

	_, err = w.Write(buf.Bytes())
	return err
}

func (te *TemplateEngine) ExecuteError(w http.ResponseWriter, r *http.Request, baseUrl string, statuscode int, tmplerr error) error {
	w.WriteHeader(statuscode)
	t, err := te.GetTemplate("pages/error.html")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return err
	}

	rd := &RenderData{
		BaseURL: baseUrl,
		Data:    tmplerr,
	}
	return t.ExecuteTemplate(w, "layout", rd)
}
