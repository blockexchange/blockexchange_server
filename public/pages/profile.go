package pages

import (
	"net/http"
)

type ProfileModel struct {
}

func (ctrl *Controller) Profile(w http.ResponseWriter, r *http.Request) {
	ctrl.te.Execute("pages/profile.html", w, r, "./", nil)
}
