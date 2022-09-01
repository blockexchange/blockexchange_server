package pages

import (
	"blockexchange/controller"
	"blockexchange/db"
	"blockexchange/types"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type SchemaModel struct {
	Schema *types.SchemaSearchResult
}

func searchSchema(sr db.SchemaSearchRepository, r *http.Request) (*types.SchemaSearchResult, error) {
	vars := mux.Vars(r)
	username := vars["username"]
	schemaname := vars["schemaname"]

	list, err := sr.Search(&types.SchemaSearchRequest{
		UserName:   &username,
		SchemaName: &schemaname,
	}, 1, 0)

	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, errors.New("schema not found")
	}

	return list[0], nil
}

func Schema(rc *controller.RenderContext) error {
	r := rc.Request()
	schema, err := searchSchema(rc.Repositories().SchemaSearchRepo, r)
	if err != nil {
		return err
	}
	m := SchemaModel{
		Schema: schema,
	}

	if m.Schema == nil {
		return errors.New("not found")
	}

	return rc.Render("pages/schema.html", m)
}
