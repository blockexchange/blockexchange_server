package schema

import (
	"blockexchange/controller"
	"blockexchange/public/components"
	"blockexchange/types"
	"errors"
	"net/http"
)

type SchemaDeleteModel struct {
	Schema     *types.SchemaSearchResult
	Breadcrumb *components.BreadcrumbModel
	Confirm    string
}

func SchemaDelete(rc *controller.RenderContext) error {
	r := rc.Request()
	repos := rc.Repositories()
	m := &SchemaDeleteModel{}

	schema, err := searchSchema(repos.SchemaSearchRepo, r)
	if err != nil {
		return err
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		if r.FormValue("confirm") == "true" {
			if schema.UserID != rc.Claims().UserID && !rc.Claims().IsAdmin() {
				return errors.New("unauthorized")
			}

			err = repos.SchemaRepo.DeleteSchema(schema.ID, schema.UserID)
			if err != nil {
				return err
			}

			rc.Redirect("../../" + schema.UserName)
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

	return rc.Render("pages/schema/schema_delete.html", m)
}
