package schema

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/types"
	"blockexchange/web/components"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type SchemaModel struct {
	Schema     *types.SchemaSearchResult
	Starred    bool
	Breadcrumb *components.BreadcrumbModel
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

func (sc *SchemaContext) Schema(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	schema, err := searchSchema(sc.repos.SchemaSearchRepo, r)
	if err != nil {
		sc.tu.RenderError(w, r, 500, err)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		if r.FormValue("action") == "update-screenshot" {
			if schema.UserID != claims.UserID && !claims.HasPermission(types.JWTPermissionAdmin) {
				sc.tu.RenderError(w, r, 500, errors.New("you are not the owner of the schema"))
				return
			}

			// update screenshot
			_, err = core.UpdatePreview(&schema.Schema, &sc.repos)
			if err != nil {
				sc.tu.RenderError(w, r, 500, err)
				return
			}
		}
	}

	/*
		rc.AddMetaTag("og:title", fmt.Sprintf("%s by %s", schema.Name, schema.UserName))
		rc.AddMetaTag("og:site_name", "Block exchange")
		rc.AddMetaTag("og:type", "Schematic")
		rc.AddMetaTag("og:url", fmt.Sprintf("%s/schema/%s/%s", rc.Config().BaseURL, schema.UserName, schema.Name))
		rc.AddMetaTag("og:image", fmt.Sprintf("%s/api/schema/%d/screenshot", rc.Config().BaseURL, schema.ID))
		if schema.Description != "" {
			rc.AddMetaTag("og:description", schema.Description)
		}
	*/

	if r.Method == http.MethodPost && claims != nil {
		r.ParseForm()
		switch r.FormValue("action") {
		case "star":
			err = sc.repos.SchemaStarRepo.Create(schema.ID, claims.UserID)
		case "unstar":
			err = sc.repos.SchemaStarRepo.Delete(schema.ID, claims.UserID)
		}
		if err != nil {
			sc.tu.RenderError(w, r, 500, err)
			return
		}

		// refresh schema
		schema, err = searchSchema(sc.repos.SchemaSearchRepo, r)
		if err != nil {
			sc.tu.RenderError(w, r, 500, err)
			return
		}
	}

	err = sc.repos.SchemaRepo.IncrementViews(schema.ID)
	if err != nil {
		sc.tu.RenderError(w, r, 500, err)
		return
	}

	m := SchemaModel{
		Schema: schema,
		Breadcrumb: components.Breadcrumb(
			components.BreadcrumbEntry{Name: "Home", Link: "/"},
			components.BreadcrumbEntry{Name: "Users", Link: "/users"},
			components.BreadcrumbEntry{Name: schema.UserName, Link: "/schema/" + schema.UserName},
			components.BreadcrumbEntry{Name: schema.Name},
		),
	}

	if claims != nil {
		star, err := sc.repos.SchemaStarRepo.GetBySchemaAndUserID(schema.ID, claims.UserID)
		if err != nil {
			sc.tu.RenderError(w, r, 500, err)
			return
		}
		m.Starred = star != nil
	}

	if m.Schema == nil {
		sc.tu.RenderError(w, r, 404, errors.New("not found"))
		return
	}

	sc.tu.ExecuteTemplate(w, r, "schema/schema.html", m)
}
