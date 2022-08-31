package pages

import (
	"blockexchange/controller"
	"net/http"
)

func About(rc *controller.RenderContext, r *http.Request) error {
	return rc.Render("pages/about.html", nil)
}
