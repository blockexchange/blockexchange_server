package web

import (
	"fmt"
	"net/http"
)

func (ctx *Context) NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx.tu.RenderError(w, r, http.StatusNotFound, fmt.Errorf("not found: %s", r.URL.Path))
	}
}
