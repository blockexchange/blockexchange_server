package pages

import (
	"net/http"
)

type SearchModel struct {
}

func (ctrl *Controller) Search(w http.ResponseWriter, r *http.Request) {
	ctrl.te.Execute("pages/search.html", w, r, "./", nil)
}
