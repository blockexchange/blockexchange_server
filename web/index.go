package web

import (
	"blockexchange/web/components"
	"net/http"
)

func (ctx *Context) Index() http.HandlerFunc {
	t := ctx.CreateTemplate("index.html")

	return func(w http.ResponseWriter, r *http.Request) {
		LatestSchemas, err := components.LatestSchemas(ctx.Repos)
		if err != nil {
			ctx.RenderError(w, r, 500, err)
			return
		}
		t.ExecuteTemplate(w, "layout", map[string]any{
			"LatestSchemas": LatestSchemas,
		})
	}
}
