package core

import (
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

func (c *Core) UpdatePreview(schema *types.Schema) (*types.SchemaScreenshot, error) {

	renderer := render.NewISORenderer(c.repos.SchemaPartRepo.GetBySchemaUIDAndOffset, cm)
	png, err := renderer.RenderIsometricPreview(schema)
	if err != nil {
		return nil, err
	}

	screenshots, err := c.repos.SchemaScreenshotRepo.GetBySchemaUID(schema.UID)
	if err != nil {
		return nil, err
	}

	var screenshot *types.SchemaScreenshot

	if len(screenshots) >= 1 {
		// update existing
		screenshot = screenshots[0]
		screenshot.Data = png

		err = c.repos.SchemaScreenshotRepo.Update(screenshot)
		if err != nil {
			return nil, err
		}
	} else {
		// create a new one
		screenshot = &types.SchemaScreenshot{
			SchemaUID: schema.UID,
			Type:      "image/png",
			Title:     "Isometric preview",
			Data:      png,
		}

		err = c.repos.SchemaScreenshotRepo.Create(screenshot)
		if err != nil {
			return nil, err
		}
	}

	return screenshot, nil
}
