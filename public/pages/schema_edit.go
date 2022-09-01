package pages

import (
	"blockexchange/controller"
	"blockexchange/types"
	"errors"
)

type SchemaEditModel struct {
	Schema *types.SchemaSearchResult
}

func SchemaEdit(rc *controller.RenderContext) error {
	r := rc.Request()
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
