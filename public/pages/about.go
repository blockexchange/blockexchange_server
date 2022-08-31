package pages

import (
	"net/http"
)

func About(rc *RenderContext, r *http.Request) error {
	return rc.Render("pages/about.html", nil)
}
