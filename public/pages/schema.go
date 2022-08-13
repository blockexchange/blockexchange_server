package pages

import (
	"blockexchange/types"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type SchemaModel struct {
	Schema *types.SchemaSearchResult
}

func (ctrl *Controller) Schema(w http.ResponseWriter, r *http.Request) {
	baseUrl := "../../"

	vars := mux.Vars(r)
	username := vars["username"]
	schemaname := vars["schemaname"]

	list, err := ctrl.SchemaSearchRepo.Search(&types.SchemaSearchRequest{
		UserName:   &username,
		SchemaName: &schemaname,
	}, 1, 0)

	if err != nil {
		ctrl.te.ExecuteError(w, r, baseUrl, 500, err)
		return
	}

	if len(list) == 0 {
		ctrl.te.ExecuteError(w, r, baseUrl, 400, errors.New("schema not found"))
		return
	}

	m := SchemaModel{
		Schema: list[0],
	}

	ctrl.te.Execute("pages/schema.html", w, r, baseUrl, m)
}
