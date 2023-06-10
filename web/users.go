package web

import (
	"blockexchange/types"
	"blockexchange/web/components"
	"net/http"
)

type UsersModel struct {
	Users      []*types.User
	Pager      *components.PagerModel
	Breadcrumb *components.BreadcrumbModel
}

func (ctx *Context) Users(w http.ResponseWriter, r *http.Request) {
	m := &UsersModel{}

	count, err := ctx.Repos.UserRepo.CountUsers()
	if err != nil {
		ctx.RenderError(w, r, 500, err)
		return
	}

	m.Pager = components.Pager(r, 20, count)

	m.Users, err = ctx.Repos.UserRepo.GetUsers(20, m.Pager.Offset)
	if err != nil {
		ctx.RenderError(w, r, 500, err)
		return
	}

	m.Breadcrumb = components.Breadcrumb(
		components.BreadcrumbEntry{Name: "Home", Link: "/"},
		components.BreadcrumbEntry{Name: "Users"},
	)

	ctx.ExecuteTemplate(w, r, "users.html", m)
}
