package db

import (
	"blockexchange/types"
	"context"

	"github.com/vingarcia/ksql"
)

var schemamodTable = ksql.NewTable("schemamod", "schema_uid", "mod_name")

type SchemaModRepository struct {
	kdb ksql.Provider
}

func (r *SchemaModRepository) GetSchemaModsBySchemaUID(schema_uid string) ([]*types.SchemaMod, error) {
	list := []*types.SchemaMod{}
	return list, r.kdb.Query(context.Background(), &list, "from schemamod where schema_uid = $1", schema_uid)
}

func (r *SchemaModRepository) GetSchemaModsBySchemaUIDs(schema_uids []string) ([]*types.SchemaMod, error) {
	list := []*types.SchemaMod{}
	return list, r.kdb.Query(context.Background(), &list, "from schemamod where schema_uid = any($1::uuid[])", schema_uids)
}

func (r *SchemaModRepository) CreateSchemaMod(schema_mod *types.SchemaMod) error {
	return r.kdb.Insert(context.Background(), schemamodTable, schema_mod)
}

func (r *SchemaModRepository) RemoveSchemaMods(schema_uid string) error {
	_, err := r.kdb.Exec(context.Background(), "delete from schemamod where schema_uid = $1", schema_uid)
	return err
}
