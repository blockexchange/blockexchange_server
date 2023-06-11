package schema

import (
	"blockexchange/types"
	"blockexchange/web/components"
	"net/http"
)

func (sc *SchemaContext) UserSchema(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	sr := sc.repos.SchemaSearchRepo

	username, _ := extractUsernameSchema(r)

	complete := true
	q := &types.SchemaSearchRequest{UserName: &username, Complete: &complete}
	count, err := sr.Count(q)
	if err != nil {
		sc.tu.RenderError(w, r, 500, err)
		return
	}
	pager := components.Pager(r, 20, count)

	list, err := sr.Search(q, 20, pager.Offset)
	if err != nil {
		sc.tu.RenderError(w, r, 500, err)
		return
	}

	m := make(map[string]any)
	m["Username"] = username
	m["Pager"] = pager
	m["SchemaList"] = components.SchemaList(c, list, false)
	m["Breadcrumb"] = components.Breadcrumb(
		components.BreadcrumbEntry{Name: "Home", Link: "/"},
		components.BreadcrumbEntry{Name: "Users", Link: "/users"},
		components.BreadcrumbEntry{Name: username},
	)

	sc.tu.ExecuteTemplate(w, r, "schema/user_schemas.html", m)
}
