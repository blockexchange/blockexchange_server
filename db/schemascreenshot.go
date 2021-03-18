package db

import "blockexchange/types"

type SchemaScreenshotRepository interface {
	GetBySchemaID(schema_id int64) ([]types.SchemaScreenshot, error)
	GetByID(id int64) (*types.SchemaScreenshot, error)
	Create(screenshot *types.SchemaScreenshot) error
	Delete(id int64) error
}

type DBSchemaScreenshotRepository struct {
}

func (r *DBSchemaScreenshotRepository) GetBySchemaID(schema_id int64) ([]types.SchemaScreenshot, error) {
	return nil, nil
}
