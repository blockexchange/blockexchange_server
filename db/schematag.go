package db

import (
	"blockexchange/types"

	"github.com/jmoiron/sqlx"
)

type SchemaTagRepository interface {
	Create(schema_id, tag_id int64) error
	Delete(schema_id, tag_id int64) error
	GetBySchemaID(schema_id int64) ([]types.SchemaTag, error)
	GetByTagID(tag_id int64) ([]types.SchemaTag, error)
}

type DBSchemaTagRepository struct {
	DB *sqlx.DB
}

func (repo DBSchemaTagRepository) Create(schema_id, tag_id int64) error {
	_, err := repo.DB.Exec("insert into schematag(schema_id, tag_id) values($1, $2)", schema_id, tag_id)
	return err
}

func (repo DBSchemaTagRepository) Delete(schema_id, tag_id int64) error {
	_, err := repo.DB.Exec("delete from schematag where schema_id = $1 and tag_id = $2", schema_id, tag_id)
	return err
}

func (repo DBSchemaTagRepository) GetBySchemaID(schema_id int64) ([]types.SchemaTag, error) {
	list := []types.SchemaTag{}
	query := `select * from schematag where schema_id = $1`
	err := repo.DB.Select(&list, query, schema_id)
	if err != nil {
		return nil, err
	} else {
		return list, nil
	}
}

func (repo DBSchemaTagRepository) GetByTagID(tag_id int64) ([]types.SchemaTag, error) {
	list := []types.SchemaTag{}
	query := `select * from schematag where tag_id = $1`
	err := repo.DB.Select(&list, query, tag_id)
	if err != nil {
		return nil, err
	} else {
		return list, nil
	}
}
