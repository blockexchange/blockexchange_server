package schema

import (
	"blockexchange/controller"
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/public/components"
	"blockexchange/types"
	"errors"
	"fmt"
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

func Schema(rc *controller.RenderContext) error {
	r := rc.Request()
	repos := rc.Repositories()
	claims := rc.Claims()

	schema, err := searchSchema(repos.SchemaSearchRepo, r)
	if err != nil {
		return err
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		if r.FormValue("action") == "update-screenshot" {
			if schema.UserID != claims.UserID && !claims.HasPermission(types.JWTPermissionAdmin) {
				return errors.New("you are not the owner of the schema")
			}

			// update screenshot
			_, err = core.UpdatePreview(&schema.Schema, repos)
			if err != nil {
				return err
			}
		}
	}

	rc.AddMetaTag("og:title", fmt.Sprintf("%s by %s", schema.Name, schema.UserName))
	rc.AddMetaTag("og:site_name", "Block exchange")
	rc.AddMetaTag("og:type", "Schematic")
	rc.AddMetaTag("og:url", fmt.Sprintf("%s/schema/%s/%s", rc.Config().BaseURL, schema.UserName, schema.Name))
	rc.AddMetaTag("og:image", fmt.Sprintf("%s/api/schema/%d/screenshot", rc.Config().BaseURL, schema.ID))
	if schema.Description != "" {
		rc.AddMetaTag("og:description", schema.Description)
	}

	if r.Method == http.MethodPost && claims != nil {
		r.ParseForm()
		switch r.FormValue("action") {
		case "star":
			err = repos.SchemaStarRepo.Create(schema.ID, claims.UserID)
		case "unstar":
			err = repos.SchemaStarRepo.Delete(schema.ID, claims.UserID)
		}
		if err != nil {
			return err
		}

		// refresh schema
		schema, err = searchSchema(repos.SchemaSearchRepo, r)
		if err != nil {
			return err
		}
	}

	err = repos.SchemaRepo.IncrementViews(schema.ID)
	if err != nil {
		return err
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
		star, err := repos.SchemaStarRepo.GetBySchemaAndUserID(schema.ID, claims.UserID)
		if err != nil {
			return err
		}
		m.Starred = star != nil
	}

	if m.Schema == nil {
		return errors.New("not found")
	}

	return rc.Render("pages/schema/schema.html", m)
}
