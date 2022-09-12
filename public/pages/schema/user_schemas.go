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

	q := &types.SchemaSearchRequest{UserName: &username}
	count, err := sr.Count(q)
	if err != nil {
		return err
	}

	list, err := sr.Search(q, 20, 0)
	if err != nil {
		return err
	}

	m := make(map[string]any)
	m["Username"] = username
	m["SchemaList"] = components.SchemaList(rc, list, false)
	m["Pager"] = components.Pager(rc, 20, count)
	m["Breadcrumb"] = components.Breadcrumb(
		components.BreadcrumbEntry{Name: "Home", Link: "/"},
		components.BreadcrumbEntry{Name: "Users", Link: "/users"},
		components.BreadcrumbEntry{Name: username},
	)

	return rc.Render("pages/schema/user_schemas.html", m)
}
