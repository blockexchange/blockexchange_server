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
	Pager      *components.PagerModel
	Query      string
	TagID      int64
	Offset     int
	Limit      int
}

func Search(rc *controller.RenderContext) error {
	repos := rc.Repositories()
	tags, err := repos.TagRepo.GetAll()
	if err != nil {
		return err
	}

	m := &SearchModel{
		Tags:  tags,
		Query: rc.Request().URL.Query().Get("q"),
		Limit: 20,
	}

	tagidstr := rc.Request().URL.Query().Get("tagid")
	if tagidstr != "" {
		tagid, err := strconv.ParseInt(tagidstr, 10, 64)
		if err != nil {
			return err
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

	count, err := repos.SchemaSearchRepo.Count(q)
	if err != nil {
		return err
	}

	m.Pager = components.Pager(rc, m.Limit, count)

	list, err := repos.SchemaSearchRepo.Search(q, m.Limit, m.Pager.Offset)
	if err != nil {
		return err
	}
	m.SchemaList = components.SchemaList(rc, list)

	return rc.Render("pages/search.html", m)
}
