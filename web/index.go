package web

import (
	"blockexchange/web/components"
	"net/http"
)

func (ctx *Context) Index(w http.ResponseWriter, r *http.Request) {
	LatestSchemas, err := components.LatestSchemas(ctx.Repos)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	ctx.tu.ExecuteTemplate(w, r, "index.html", map[string]any{
		"LatestSchemas": LatestSchemas,
	})
}
