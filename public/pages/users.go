package pages

import (
	"blockexchange/controller"
)

type UsersModel struct {
}

func Users(rc *controller.RenderContext) error {
	return rc.Render("pages/users.html", nil)
}
