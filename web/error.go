package web

import (
	"fmt"
	"net/http"
)

func (ctx *Context) RenderError(w http.ResponseWriter, r *http.Request, code int, err error) {
	w.WriteHeader(code)
	err2 := ctx.error_template.ExecuteTemplate(w, "layout", err)
	if err2 != nil {
		//https://stackoverflow.com/questions/44675087/golang-template-variable-isset
		panic(err2)
	}
}

func (ctx *Context) NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx.RenderError(w, r, http.StatusNotFound, fmt.Errorf("not found: %s", r.URL.Path))
	}
}
