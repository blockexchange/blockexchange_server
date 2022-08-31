package pages

import (
	"blockexchange/controller"
	"blockexchange/types"
	"net/http"
)

type TagsModel struct {
}

func Tags(rc *controller.RenderContext, r *http.Request, c *types.Claims) error {
	return rc.Render("pages/tags.html", nil)
}
