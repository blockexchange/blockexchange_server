package pages

import (
	"net/http"
)

func (ctrl *Controller) About(w http.ResponseWriter, r *http.Request) {
	ctrl.te.Execute("pages/about.html", w, r, "./", nil)
}
