package controller

import "net/http"

func (ctrl *Controller) Index(w http.ResponseWriter, r *http.Request) {
	ctrl.te.Execute("pages/index.html", w, r, "./", nil)
}
