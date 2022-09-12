package schema

import (
	"blockexchange/controller"
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

	schema, err := searchSchema(repos.SchemaSearchRepo, r)
	if err != nil {
		return err
	}

	rc.AddMetaTag("og:title", fmt.Sprintf("%s by %s", schema.Name, schema.UserName))
	rc.AddMetaTag("og:site_name", "Block exchange")
	rc.AddMetaTag("og:type", "Schematic")
	rc.AddMetaTag("og:url", fmt.Sprintf("%s/schema/%s/%s", rc.Config().BaseURL, schema.UserName, schema.Name))
	rc.AddMetaTag("og:image", fmt.Sprintf("%s/api/schema/%d/screenshot", rc.Config().BaseURL, schema.ID))
	if schema.Description != "" {
		rc.AddMetaTag("og:description", schema.Description)
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

	if m.Schema == nil {
		return errors.New("not found")
	}

	return rc.Render("pages/schema/schema.html", m)
}
