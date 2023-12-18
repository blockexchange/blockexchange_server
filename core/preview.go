package core

import (
	"blockexchange/db"
	"blockexchange/render"
	"blockexchange/types"

	"github.com/minetest-go/colormapping"
)

var cm *colormapping.ColorMapping = colormapping.NewColorMapping()

func init() {
	err := cm.LoadDefaults()
	if err != nil {
		panic(err)
	}
}

// TODO: (c *Core)
func UpdatePreview(schema *types.Schema, repo *db.Repositories) (*types.SchemaScreenshot, error) {

	renderer := render.NewISORenderer(repo.SchemaPartRepo.GetBySchemaIDAndOffset, cm)
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
