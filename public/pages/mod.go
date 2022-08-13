package pages

import (
	"net/http"
)

func (ctrl *Controller) Mod(w http.ResponseWriter, r *http.Request) {
	ctrl.te.Execute("pages/mod.html", w, r, "./", nil)
}
