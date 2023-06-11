package web

import (
	"blockexchange/types"
	"blockexchange/web/components"
	"net/http"
	"strconv"
)

type SearchModel struct {
	Tags       []*types.Tag
	SchemaList *components.SchemaListModel
	Pager      *components.PagerModel
	Breadcrumb *components.BreadcrumbModel
	Query      string
	TagID      int64
	Offset     int
	Limit      int
}

func (ctx *Context) Search(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	tags, err := ctx.Repos.TagRepo.GetAll()
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	m := &SearchModel{
		Tags:  tags,
		Query: r.URL.Query().Get("q"),
		Limit: 20,
	}

	tagidstr := r.URL.Query().Get("tagid")
	if tagidstr != "" {
		tagid, err := strconv.ParseInt(tagidstr, 10, 64)
		if err != nil {
			ctx.tu.RenderError(w, r, 500, err)
			return
		}

		m.TagID = tagid
	}

	complete := true
	q := &types.SchemaSearchRequest{Complete: &complete}
	if m.Query != "" {
		q.Keywords = &m.Query
	}
	if m.TagID > 0 {
		q.TagID = &m.TagID
	}

	count, err := ctx.Repos.SchemaSearchRepo.Count(q)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}

	m.Pager = components.Pager(r, m.Limit, count)

	list, err := ctx.Repos.SchemaSearchRepo.Search(q, m.Limit, m.Pager.Offset)
	if err != nil {
		ctx.tu.RenderError(w, r, 500, err)
		return
	}
	m.SchemaList = components.SchemaList(c, list, true)

	m.Breadcrumb = components.Breadcrumb(
		components.BreadcrumbEntry{Name: "Home", Link: "/"},
		components.BreadcrumbEntry{Name: "Search"},
	)

	ctx.tu.ExecuteTemplate(w, r, "search.html", m)
}
