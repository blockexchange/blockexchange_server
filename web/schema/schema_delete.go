package schema

import (
	"blockexchange/types"
	"blockexchange/web/components"
	"errors"
	"net/http"
)

type SchemaDeleteModel struct {
	Schema     *types.SchemaSearchResult
	Breadcrumb *components.BreadcrumbModel
	Confirm    string
}

func (sc *SchemaContext) SchemaDelete(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	m := &SchemaDeleteModel{}

	schema, err := searchSchema(sc.repos.SchemaSearchRepo, r)
	if err != nil {
		sc.tu.RenderError(w, r, 500, err)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		if r.FormValue("confirm") == "true" {
			if schema.UserID != claims.UserID && !claims.IsAdmin() {
				sc.tu.RenderError(w, r, 401, errors.New("unauthorized"))
				return
			}

			err = sc.repos.SchemaRepo.DeleteSchema(schema.ID, schema.UserID)
			if err != nil {
				sc.tu.RenderError(w, r, 500, err)
				return
			}

			http.Redirect(w, r, sc.BaseURL+"/schema/"+schema.UserName, http.StatusSeeOther)
		}
	}

	m.Schema = schema
	m.Breadcrumb = components.Breadcrumb(
		components.BreadcrumbEntry{Name: "Home", Link: "/"},
		components.BreadcrumbEntry{Name: "Users", Link: "/users"},
		components.BreadcrumbEntry{Name: schema.UserName, Link: "/schema/" + schema.UserName},
		components.BreadcrumbEntry{Name: schema.Name, Link: "/schema/" + schema.UserName + "/" + schema.Name},
		components.BreadcrumbEntry{Name: "Delete"},
	)

	sc.tu.ExecuteTemplate(w, r, "schema/schema_delete.html", m)
}
