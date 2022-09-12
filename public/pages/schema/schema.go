package schema

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

func extractUsernameSchema(r *http.Request) (string, string) {
	vars := mux.Vars(r)
	return vars["username"], vars["schemaname"]
}

func searchSchema(sr db.SchemaSearchRepository, r *http.Request) (*types.SchemaSearchResult, error) {
	username, schemaname := extractUsernameSchema(r)

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

	return rc.Render("pages/schema/schema.html", m)
}
