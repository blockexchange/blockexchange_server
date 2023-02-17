package templateengine

import (
	"bytes"
	"html/template"
	"io/fs"
	"net/http"
	"sync"
)

type TemplateEngineOptions struct {
	Templates   fs.FS
	EnableCache bool
	FuncMap     map[string]any
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

	tmpl, err := template.New("").Funcs(te.options.FuncMap).Option("missingkey=error").ParseFS(f, "layout.html", "components/*.html")
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

func (te *TemplateEngine) Execute(file string, w http.ResponseWriter, r *http.Request, statuscode int, data any) error {
	t, err := te.GetTemplate(file)
	if err != nil {
		return err
	}

	buf := bytes.Buffer{}

	err = t.ExecuteTemplate(&buf, "layout", data)
	if err != nil {
		return err
	}

	_, err = w.Write(buf.Bytes())
	return err
}
