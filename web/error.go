package web

import (
	"fmt"
	"net/http"
)

var error_template = createTemplate("error.html")

func RenderError(w http.ResponseWriter, r *http.Request, code int, err error) {
	w.WriteHeader(code)
	error_template.ExecuteTemplate(w, "layout", err)
}

func (ctx *Context) NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RenderError(w, r, http.StatusNotFound, fmt.Errorf("not found: %s", r.URL.Path))
	}
}
