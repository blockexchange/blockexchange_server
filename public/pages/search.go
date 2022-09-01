package pages

import (
	"blockexchange/controller"
)

type SearchModel struct {
}

func Search(rc *controller.RenderContext) error {
	return rc.Render("pages/search.html", nil)
}
