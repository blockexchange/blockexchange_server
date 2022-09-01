package pages

import (
	"blockexchange/controller"
)

func Mod(rc *controller.RenderContext) error {
	return rc.Render("pages/mod.html", nil)
}
