package pages

import (
	"blockexchange/controller"
	"blockexchange/public/components"
)

type IndexModel struct {
	LatestSchemas *components.LatestSchemasModel
}

func Index(rc *controller.RenderContext) error {
	m := &IndexModel{}

	var err error
	m.LatestSchemas, err = components.LatestSchemas(rc.BaseURL(), rc.Repositories())
	if err != nil {
		return err
	}

	return rc.Render("pages/index.html", m)
}
