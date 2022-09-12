package pages

import (
	"blockexchange/controller"
	"blockexchange/public/components"
	"blockexchange/types"
)

type UsersModel struct {
	Users      []*types.User
	Pager      *components.PagerModel
	Breadcrumb *components.BreadcrumbModel
}

func Users(rc *controller.RenderContext) error {
	m := &UsersModel{}

	count, err := rc.Repositories().UserRepo.CountUsers()
	if err != nil {
		return err
	}

	m.Pager = components.Pager(rc, 20, count)

	m.Users, err = rc.Repositories().UserRepo.GetUsers(20, m.Pager.Offset)
	if err != nil {
		return err
	}

	m.Breadcrumb = components.Breadcrumb(
		components.BreadcrumbEntry{Name: "Home", Link: "/"},
		components.BreadcrumbEntry{Name: "Users"},
	)

	return rc.Render("pages/users.html", m)
}
