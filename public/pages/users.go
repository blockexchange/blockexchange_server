package pages

import (
	"blockexchange/controller"
	"net/http"
)

type UsersModel struct {
}

func Users(rc *controller.RenderContext, r *http.Request) error {
	return rc.Render("pages/users.html", nil)
}
