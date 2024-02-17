package db

import (
	"blockexchange/types"
	"context"

	"github.com/google/uuid"
	"github.com/vingarcia/ksql"
)

var schemaScreenshotTable = ksql.NewTable("schema_screenshot", "uid")

type SchemaScreenshotRepository struct {
	kdb ksql.Provider
}

func (r *SchemaScreenshotRepository) GetByUID(uid string) (*types.SchemaScreenshot, error) {
	result := &types.SchemaScreenshot{}
	err := r.kdb.QueryOne(context.Background(), result, "from schema_screenshot where uid = $1", uid)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	} else {
		return result, err
	}
}

func (r *SchemaScreenshotRepository) GetBySchemaUID(schema_uid string) ([]*types.SchemaScreenshot, error) {
	list := []*types.SchemaScreenshot{}
	return list, r.kdb.Query(context.Background(), &list, "from schema_screenshot where schema_uid = $1", schema_uid)
}

func (r *SchemaScreenshotRepository) Create(screenshot *types.SchemaScreenshot) error {
	if screenshot.UID == "" {
		screenshot.UID = uuid.NewString()
	}
	return r.kdb.Insert(context.Background(), schemaScreenshotTable, screenshot)
}

func (r *SchemaScreenshotRepository) Update(screenshot *types.SchemaScreenshot) error {
	return r.kdb.Patch(context.Background(), schemaScreenshotTable, screenshot)
}
