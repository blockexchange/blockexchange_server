package db

import (
	"blockexchange/types"

	"github.com/jmoiron/sqlx"
)

type SchemaTagRepository struct {
	DB *sqlx.DB
}

func (repo SchemaTagRepository) Create(schema_id, tag_id int64) error {
	_, err := repo.DB.Exec("insert into schematag(schema_id, tag_id) values($1, $2)", schema_id, tag_id)
	return err
}

func (repo SchemaTagRepository) Delete(schema_id, tag_id int64) error {
	_, err := repo.DB.Exec("delete from schematag where schema_id = $1 and tag_id = $2", schema_id, tag_id)
	return err
}

func (repo SchemaTagRepository) GetBySchemaID(schema_id int64) ([]types.SchemaTag, error) {
	list := []types.SchemaTag{}
	query := `select * from schematag where schema_id = $1`
	err := repo.DB.Select(&list, query, schema_id)
	if err != nil {
		return nil, err
	} else {
		return list, nil
	}
}

func (repo SchemaTagRepository) GetByTagID(tag_id int64) ([]types.SchemaTag, error) {
	list := []types.SchemaTag{}
	query := `select * from schematag where tag_id = $1`
	err := repo.DB.Select(&list, query, tag_id)
	if err != nil {
		return nil, err
	} else {
		return list, nil
	}
}
