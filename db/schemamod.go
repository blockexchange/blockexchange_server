package db

import (
	"blockexchange/types"
	"context"

	"github.com/vingarcia/ksql"
)

var schemamodTable = ksql.NewTable("schemamod", "id")

type SchemaModRepository struct {
	kdb ksql.Provider
}

func (r *SchemaModRepository) GetSchemaModsBySchemaID(schema_id int64) ([]*types.SchemaMod, error) {
	list := []*types.SchemaMod{}
	return list, r.kdb.Query(context.Background(), &list, "from schemamod where schema_id = $1", schema_id)
}

func (r *SchemaModRepository) CreateSchemaMod(schema_mod *types.SchemaMod) error {
	return r.kdb.Insert(context.Background(), schemamodTable, schema_mod)
}

func (r *SchemaModRepository) RemoveSchemaMods(schema_id int64) error {
	_, err := r.kdb.Exec(context.Background(), "delete from schemamod where schema_id = $1", schema_id)
	return err
}
