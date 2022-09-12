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

	list, err := sr.Search(&types.SchemaSearchRequest{
		UserName: &username,
	}, 20, 0)

	if err != nil {
		return err
	}

	m := make(map[string]any)
	m["Username"] = username
	m["SchemaList"] = components.SchemaList(rc, list)

	return rc.Render("pages/schema/user_schemas.html", m)
}
