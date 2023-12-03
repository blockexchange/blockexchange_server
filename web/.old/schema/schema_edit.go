package schema

import (
	"blockexchange/types"
	"blockexchange/web/components"
	"errors"
	"net/http"
)

type SchemaEditModel struct {
	Schema     *types.SchemaSearchResult
	Breadcrumb *components.BreadcrumbModel
}

func (sc *SchemaContext) schemaEditPost(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	sr := sc.repos.SchemaRepo

	err := r.ParseForm()
	if err != nil {
		sc.tu.RenderError(w, r, 500, err)
		return
	}

	username, schemaname := extractUsernameSchema(r)
	schema, err := sr.GetSchemaByUsernameAndName(username, schemaname)
	if err != nil {
		sc.tu.RenderError(w, r, 500, err)
		return
	}
	if schema == nil {
		sc.tu.RenderError(w, r, 404, errors.New("not found"))
		return
	}
	if claims.UserID != schema.UserID {
		sc.tu.RenderError(w, r, 403, errors.New("not allowed"))
		return
	}

	newSchemaName := r.FormValue("name")
	newSchema, err := sr.GetSchemaByUsernameAndName(username, newSchemaName)
	if err != nil {
		sc.tu.RenderError(w, r, 403, err)
		return
	}
	if newSchema != nil && newSchema.ID != schema.ID {
		sc.tu.RenderError(w, r, 403, errors.New("schema name already taken"))
		return
	}

	schema.Description = r.FormValue("description")
	schema.License = r.FormValue("license")
	schema.Name = newSchemaName

	err = sr.UpdateSchema(schema)
	if err != nil {
		sc.tu.RenderError(w, r, 500, err)
		return
	}

	http.Redirect(w, r, sc.BaseURL+"/schema/"+username+"/"+schema.Name, http.StatusSeeOther)
}

func (sc *SchemaContext) SchemaEdit(w http.ResponseWriter, r *http.Request, claims *types.Claims) {

	if r.Method == http.MethodPost {
		sc.schemaEditPost(w, r, claims)
		return
	}

	schema, err := searchSchema(sc.repos.SchemaSearchRepo, r)
	if err != nil {
		sc.tu.RenderError(w, r, 500, err)
		return
	}
	if schema == nil {
		sc.tu.RenderError(w, r, 404, errors.New("not found"))
		return
	}

	m := SchemaEditModel{
		Schema: schema,
	}

	m.Breadcrumb = components.Breadcrumb(
		components.BreadcrumbEntry{Name: "Home", Link: "/"},
		components.BreadcrumbEntry{Name: "Users", Link: "/users"},
		components.BreadcrumbEntry{Name: schema.UserName, Link: "/schema/" + schema.UserName},
		components.BreadcrumbEntry{Name: schema.Name, Link: "/schema/" + schema.UserName + "/" + schema.Name},
		components.BreadcrumbEntry{Name: "Edit"},
	)

	sc.tu.ExecuteTemplate(w, r, "schema/schema_edit.html", m)
}
