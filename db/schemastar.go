package db

import (
	"blockexchange/types"

	"github.com/jmoiron/sqlx"
)

type SchemaStarRepository interface {
	Create(schema_id, user_id int64) error
	Delete(schema_id, user_id int64) error
	GetBySchemaID(schema_id int64) ([]types.SchemaStar, error)
	CountBySchemaID(schema_id int64) (int, error)
}

type DBSchemaStarRepository struct {
	DB *sqlx.DB
}

func (repo DBSchemaStarRepository) Create(schema_id, user_id int64) error {
	_, err := repo.DB.Exec("insert into user_schema_star(schema_id, user_id) values($1, $2)", schema_id, user_id)
	return err
}

func (repo DBSchemaStarRepository) Delete(schema_id, user_id int64) error {
	_, err := repo.DB.Exec("delete from user_schema_star where schema_id = $1 and user_id = $2", schema_id, user_id)
	return err
}

func (repo DBSchemaStarRepository) GetBySchemaID(schema_id int64) ([]types.SchemaStar, error) {
	list := []types.SchemaStar{}
	query := `select * from user_schema_star where schema_id = $1`
	err := repo.DB.Select(&list, query, schema_id)
	if err != nil {
		return nil, err
	} else {
		return list, nil
	}
}

func (repo DBSchemaStarRepository) CountBySchemaID(schema_id int64) (int, error) {
	query := `select count(*) from user_schema_star where schema_id = $1`
	row := repo.DB.QueryRow(query, schema_id)
	count := 0
	err := row.Scan(&count)
	return count, err
}
