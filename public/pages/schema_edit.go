package pages

import (
	"blockexchange/controller"
	"blockexchange/types"
	"errors"
	"net/http"
)

type SchemaEditModel struct {
	Schema *types.SchemaSearchResult
}

func SchemaEdit(rc *controller.RenderContext, r *http.Request, claims *types.Claims) error {
	schema, err := searchSchema(rc.Repositories().SchemaSearchRepo, r)
	if err != nil {
		return err
	}

	//TODO: edit stuff on POST
	m := SchemaEditModel{
		Schema: schema,
	}

	if m.Schema == nil {
		return errors.New("not found")
	}

	return rc.Render("pages/schema_edit.html", m)
}
