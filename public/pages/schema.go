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

func (ctrl *Controller) searchSchema(w http.ResponseWriter, r *http.Request, baseUrl string) *types.SchemaSearchResult {
	vars := mux.Vars(r)
	username := vars["username"]
	schemaname := vars["schemaname"]

	list, err := ctrl.SchemaSearchRepo.Search(&types.SchemaSearchRequest{
		UserName:   &username,
		SchemaName: &schemaname,
	}, 1, 0)

	if err != nil {
		ctrl.te.ExecuteError(w, r, baseUrl, 500, err)
		return nil
	}

	if len(list) == 0 {
		ctrl.te.ExecuteError(w, r, baseUrl, 400, errors.New("schema not found"))
		return nil
	}

	return list[0]
}

func (ctrl *Controller) Schema(w http.ResponseWriter, r *http.Request) {
	baseUrl := "../../"

	m := SchemaModel{
		Schema: ctrl.searchSchema(w, r, baseUrl),
	}

	if m.Schema == nil {
		return
	}
	ctrl.te.Execute("pages/schema.html", w, r, baseUrl, m)
}
