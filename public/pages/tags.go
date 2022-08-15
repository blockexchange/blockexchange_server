package pages

import (
	"blockexchange/types"
	"net/http"
)

type TagsModel struct {
}

func (ctrl *Controller) Tags(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	ctrl.te.Execute("pages/tags.html", w, r, "./", nil)
}
