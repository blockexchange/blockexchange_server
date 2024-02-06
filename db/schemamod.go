package db

import (
	"blockexchange/types"
	"database/sql"

	"github.com/minetest-go/dbutil"
)

func NewSchemaModRepository(DB *sql.DB) *SchemaModRepository {
	return &SchemaModRepository{
		DB:  DB,
		dbu: dbutil.New(DB, dbutil.DialectPostgres, func() *types.SchemaMod { return &types.SchemaMod{} }),
	}
}

type SchemaModRepository struct {
	DB  *sql.DB
	dbu *dbutil.DBUtil[*types.SchemaMod]
}

func (r *SchemaModRepository) GetSchemaModsBySchemaID(schema_id int64) ([]*types.SchemaMod, error) {
	return r.dbu.SelectMulti("where schema_id = %s", schema_id)
}

func (r *SchemaModRepository) CreateSchemaMod(schema_mod *types.SchemaMod) error {
	return r.dbu.InsertReturning(schema_mod, "id", &schema_mod.ID)
}

func (r *SchemaModRepository) RemoveSchemaMods(schema_id int64) error {
	return r.dbu.Delete("where schema_id = %s", schema_id)
}
