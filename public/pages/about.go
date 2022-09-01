package pages

import (
	"blockexchange/controller"
)

func About(rc *controller.RenderContext) error {
	return rc.Render("pages/about.html", nil)
}
