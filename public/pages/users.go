package pages

import (
	"net/http"
)

type UsersModel struct {
}

func (ctrl *Controller) Users(w http.ResponseWriter, r *http.Request) {
	ctrl.te.Execute("pages/users.html", w, r, "./", nil)
}
