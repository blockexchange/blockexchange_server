package core

import (
	"blockexchange/db"
	"blockexchange/render"
	"blockexchange/types"
)

func UpdatePreview(schema *types.Schema, repo *db.Repositories) (*types.SchemaScreenshot, error) {

	cm, err := render.GetColorMapping()
	if err != nil {
		return nil, err
	}

	renderer := render.NewRenderer(repo.SchemaPartRepo.GetBySchemaIDAndOffset, cm)
	png, err := renderer.RenderIsometricPreview(schema)
	if err != nil {
		return nil, err
	}

	screenshots, err := repo.SchemaScreenshotRepo.GetBySchemaID(schema.ID)
	if err != nil {
		return nil, err
	}

	var screenshot types.SchemaScreenshot

	if len(screenshots) >= 1 {
		// update existing
		screenshot = screenshots[0]
		screenshot.Data = png

		err = repo.SchemaScreenshotRepo.Update(&screenshot)
		if err != nil {
			return nil, err
		}
	} else {
		// create a new one
		screenshot = types.SchemaScreenshot{
			SchemaID: schema.ID,
			Type:     "image/png",
			Title:    "Isometric preview",
			Data:     png,
		}

		err = repo.SchemaScreenshotRepo.Create(&screenshot)
		if err != nil {
			return nil, err
		}
	}

	return &screenshot, nil
}
