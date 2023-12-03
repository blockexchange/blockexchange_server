package web

import (
	"blockexchange/types"
	"net/http"
)

type TagsModel struct {
	Tags []*types.Tag
}

func (ctx *Context) Tags(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	m := &TagsModel{}

	tags, err := ctx.Repos.TagRepo.GetAll()
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}
	m.Tags = tags

	ctx.tu.ExecuteTemplate(w, r, "tags.html", m)
}
