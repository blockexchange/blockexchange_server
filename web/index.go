package web

import (
	"blockexchange/types"
	"blockexchange/web/components"
	"net/http"
)

func (ctx *Context) Index(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	t := ctx.CreateTemplate("index.html")

	LatestSchemas, err := components.LatestSchemas(ctx.Repos)
	if err != nil {
		ctx.RenderError(w, r, 500, err)
		return
	}
	t.ExecuteTemplate(w, "layout", map[string]any{
		"LatestSchemas": LatestSchemas,
		"Claims":        c,
	})
}
