package controller

import (
	"blockexchange/controller/components"
	"net/http"
)

type IndexModel struct {
	LatestSchemas *components.LatestSchemasModel
}

func (ctrl *Controller) Index(w http.ResponseWriter, r *http.Request) {
	m := &IndexModel{}

	var err error
	m.LatestSchemas, err = components.LatestSchemas(ctrl.Repositories)
	if err != nil {
		ctrl.te.ExecuteError(w, r, "./", 500, err)
		return
	}

	ctrl.te.Execute("pages/index.html", w, r, "./", m)
}
