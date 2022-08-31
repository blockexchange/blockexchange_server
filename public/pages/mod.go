package pages

import (
	"blockexchange/controller"
	"net/http"
)

func Mod(rc *controller.RenderContext, r *http.Request) error {
	return rc.Render("pages/mod.html", nil)
}
