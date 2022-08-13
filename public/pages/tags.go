package pages

import (
	"net/http"
)

type TagsModel struct {
}

func (ctrl *Controller) Tags(w http.ResponseWriter, r *http.Request) {
	ctrl.te.Execute("pages/tags.html", w, r, "./", nil)
}
