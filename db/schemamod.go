package db

import (
	"blockexchange/types"

	"github.com/jmoiron/sqlx"
)

type SchemaModRepository struct {
	DB *sqlx.DB
}

func (repo SchemaModRepository) GetSchemaModsBySchemaID(schema_id int64) ([]types.SchemaMod, error) {
	list := []types.SchemaMod{}
	query := `select * from schemamod where schema_id = $1`
	err := repo.DB.Select(&list, query, schema_id)
	if err != nil {
		return nil, err
	} else {
		return list, nil
	}
}

func (repo SchemaModRepository) CreateSchemaMod(schema_mod *types.SchemaMod) error {
	query := `
		insert into
		schemamod(schema_id, mod_name)
		values(:schema_id, :mod_name)
		returning id
	`
	stmt, err := repo.DB.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&schema_mod.ID, schema_mod)
}

func (repo SchemaModRepository) RemoveSchemaMods(schema_id int64) error {
	query := `
		delete
		from schemamod
		where schema_id = $1
	`
	_, err := repo.DB.Exec(query, schema_id)
	return err
}
