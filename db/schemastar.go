package db

import (
	"blockexchange/types"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type SchemaStarRepository struct {
	DB *sqlx.DB
}

func (repo SchemaStarRepository) Create(schema_id, user_id int64) error {
	_, err := repo.DB.Exec("insert into user_schema_star(schema_id, user_id) values($1, $2)", schema_id, user_id)
	return err
}

func (repo SchemaStarRepository) Delete(schema_id, user_id int64) error {
	_, err := repo.DB.Exec("delete from user_schema_star where schema_id = $1 and user_id = $2", schema_id, user_id)
	return err
}

func (repo SchemaStarRepository) GetBySchemaID(schema_id int64) ([]types.SchemaStar, error) {
	list := []types.SchemaStar{}
	query := `select * from user_schema_star where schema_id = $1`
	err := repo.DB.Select(&list, query, schema_id)
	if err != nil {
		return nil, err
	} else {
		return list, nil
	}
}

func (repo SchemaStarRepository) GetBySchemaAndUserID(schema_id int64, user_id int64) (*types.SchemaStar, error) {
	query := `select * from user_schema_star where schema_id = $1 and user_id = $2`
	star := types.SchemaStar{}
	err := repo.DB.Get(&star, query, schema_id, user_id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &star, nil
	}
}

func (repo SchemaStarRepository) CountBySchemaID(schema_id int64) (int, error) {
	query := `select count(*) from user_schema_star where schema_id = $1`
	row := repo.DB.QueryRow(query, schema_id)
	count := 0
	err := row.Scan(&count)
	return count, err
}

func (repo SchemaStarRepository) CountByUserID(user_id int64) (int, error) {
	query := `select count(*) from user_schema_star uss join schema s on s.id = uss.schema_id where s.user_id = $1`
	row := repo.DB.QueryRow(query, user_id)
	count := 0
	err := row.Scan(&count)
	return count, err
}
