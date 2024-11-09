package db

import (
	"blockexchange/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SchemaScreenshotRepository struct {
	g *gorm.DB
}

func (r *SchemaScreenshotRepository) GetByUID(uid string) (*types.SchemaScreenshot, error) {
	return FindSingle[types.SchemaScreenshot](r.g.Where(types.SchemaScreenshot{UID: uid}))
}

func (r *SchemaScreenshotRepository) GetAllBySchemaUID(schema_uid string) ([]*types.SchemaScreenshot, error) {
	return FindMulti[types.SchemaScreenshot](r.g.Where(types.SchemaScreenshot{SchemaUID: schema_uid}))
}

func (r *SchemaScreenshotRepository) GetLatestBySchemaUID(schema_uid string) (*types.SchemaScreenshot, error) {
	return FindSingle[types.SchemaScreenshot](r.g.Where(types.SchemaScreenshot{SchemaUID: schema_uid}).Order("created desc"))
}

func (r *SchemaScreenshotRepository) Create(screenshot *types.SchemaScreenshot) error {
	if screenshot.UID == "" {
		screenshot.UID = uuid.NewString()
	}
	return r.g.Create(screenshot).Error
}

func (r *SchemaScreenshotRepository) Update(screenshot *types.SchemaScreenshot) error {
	return r.g.Updates(screenshot).Error
}
