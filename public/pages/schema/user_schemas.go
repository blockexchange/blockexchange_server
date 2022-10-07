package schema

import (
	"blockexchange/controller"
	"blockexchange/public/components"
	"blockexchange/types"
)

func UserSchema(rc *controller.RenderContext) error {
	r := rc.Request()
	sr := rc.Repositories().SchemaSearchRepo

	username, _ := extractUsernameSchema(r)

	complete := true
	q := &types.SchemaSearchRequest{UserName: &username, Complete: &complete}
	count, err := sr.Count(q)
	if err != nil {
		return err
	}
	pager := components.Pager(rc, 20, count)

	list, err := sr.Search(q, 20, pager.Offset)
	if err != nil {
		return err
	}

	m := make(map[string]any)
	m["Username"] = username
	m["Pager"] = pager
	m["SchemaList"] = components.SchemaList(rc, list, false)
	m["Breadcrumb"] = components.Breadcrumb(
		components.BreadcrumbEntry{Name: "Home", Link: "/"},
		components.BreadcrumbEntry{Name: "Users", Link: "/users"},
		components.BreadcrumbEntry{Name: username},
	)

	return rc.Render("pages/schema/user_schemas.html", m)
}
