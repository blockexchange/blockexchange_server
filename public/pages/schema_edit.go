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

func schemaEditPost(rc *controller.RenderContext) error {
	r := rc.Request()
	sr := rc.Repositories().SchemaRepo

	err := r.ParseForm()
	if err != nil {
		return err
	}

	claims := rc.Claims()
	username, schemaname := extractUsernameSchema(r)
	schema, err := sr.GetSchemaByUsernameAndName(username, schemaname)
	if err != nil {
		return err
	}
	if schema == nil {
		return errors.New("not found")
	}
	if claims.UserID != schema.UserID {
		return errors.New("not allowed")
	}

	newSchemaName := r.FormValue("name")
	newSchema, err := sr.GetSchemaByUsernameAndName(username, newSchemaName)
	if err != nil {
		return err
	}
	if newSchema != nil && newSchema.ID != schema.ID {
		return errors.New("schema name already taken")
	}

	schema.Description = r.FormValue("description")
	schema.Name = newSchemaName

	rc.Redirect("../" + schema.Name)

	return sr.UpdateSchema(schema)
}

func SchemaEdit(rc *controller.RenderContext) error {
	r := rc.Request()

	if r.Method == http.MethodPost {
		return schemaEditPost(rc)
	}

	schema, err := searchSchema(rc.Repositories().SchemaSearchRepo, r)
	if err != nil {
		return err
	}
	if schema == nil {
		return errors.New("not found")
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
