package pages

import (
	"blockexchange/controller"
	"blockexchange/types"
)

type TagsModel struct {
	Tags []*types.Tag
}

func Tags(rc *controller.RenderContext) error {
	m := &TagsModel{}

	tags, err := rc.Repositories().TagRepo.GetAll()
	if err != nil {
		return err
	}
	m.Tags = tags

	return rc.Render("pages/tags.html", m)
}
