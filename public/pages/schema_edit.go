package pages

import (
	"blockexchange/types"
	"net/http"
)

type SchemaEditModel struct {
	Schema *types.SchemaSearchResult
}

func (ctrl *Controller) SchemaEdit(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	baseUrl := "../../../"

	m := SchemaEditModel{
		Schema: ctrl.searchSchema(w, r, baseUrl),
	}

	if m.Schema == nil {
		return
	}

	ctrl.te.Execute("pages/schema_edit.html", w, r, baseUrl, m)
}
