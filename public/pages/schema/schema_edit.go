package schema

import (
	"blockexchange/controller"
	"blockexchange/public/components"
	"blockexchange/types"
	"errors"
	"net/http"
)

type SchemaEditModel struct {
	Schema     *types.SchemaSearchResult
	Breadcrumb *components.BreadcrumbModel
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
	schema.License = r.FormValue("license")
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

	m.Breadcrumb = components.Breadcrumb(
		components.BreadcrumbEntry{Name: "Home", Link: "/"},
		components.BreadcrumbEntry{Name: "Users", Link: "/users"},
		components.BreadcrumbEntry{Name: schema.UserName, Link: "/schema/" + schema.UserName},
		components.BreadcrumbEntry{Name: schema.Name, Link: "/schema/" + schema.UserName + "/" + schema.Name},
		components.BreadcrumbEntry{Name: "Edit"},
	)

	return rc.Render("pages/schema/schema_edit.html", m)
}
