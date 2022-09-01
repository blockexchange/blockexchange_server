package pages

import (
	"blockexchange/controller"
)

type TagsModel struct {
}

func Tags(rc *controller.RenderContext) error {
	return rc.Render("pages/tags.html", nil)
}
