package web

import (
	"fmt"
	"net/http"
)

func (ctx *Context) RenderError(w http.ResponseWriter, r *http.Request, code int, err error) {
	w.WriteHeader(code)
	t := ctx.CreateTemplate("error.html", r)
	t.ExecuteTemplate(w, "layout", err)
}

func (ctx *Context) NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx.RenderError(w, r, http.StatusNotFound, fmt.Errorf("not found: %s", r.URL.Path))
	}
}
