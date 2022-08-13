package pages

import (
	"net/http"
)

type SchemaModel struct {
}

func (ctrl *Controller) Schema(w http.ResponseWriter, r *http.Request) {
	ctrl.te.Execute("pages/schema.html", w, r, "../../", nil)
}
