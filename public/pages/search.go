package pages

import (
	"blockexchange/controller"
	"blockexchange/public/components"
	"blockexchange/types"
	"strconv"
)

type SearchModel struct {
	Tags       []*types.Tag
	SchemaList []*components.SchemaListEntry
	Query      string
	Offset     int
	Limit      int
}

func Search(rc *controller.RenderContext) error {
	tags, err := rc.Repositories().TagRepo.GetAll()
	if err != nil {
		return err
	}

	m := &SearchModel{
		Tags:   tags,
		Query:  rc.Request().URL.Query().Get("q"),
		Limit:  20,
		Offset: 0,
	}

	page, err := strconv.ParseInt(rc.Request().URL.Query().Get("page"), 10, 64)
	if err == nil {
		m.Offset = int(page) * m.Limit
	}

	complete := true
	q := &types.SchemaSearchRequest{Complete: &complete}
	if m.Query != "" {
		q.Keywords = &m.Query
	}
	list, err := rc.Repositories().SchemaSearchRepo.Search(q, m.Limit, m.Offset)
	if err != nil {
		return err
	}
	m.SchemaList = components.SchemaList(rc, list)

	return rc.Render("pages/search.html", m)
}
