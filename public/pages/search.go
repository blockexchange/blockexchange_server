package pages

import (
	"blockexchange/controller"
	"net/http"
)

type SearchModel struct {
}

func Search(rc *controller.RenderContext, r *http.Request) error {
	return rc.Render("pages/search.html", nil)
}
