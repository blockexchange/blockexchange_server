package db

import (
	"blockexchange/types"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type SchemaRepository interface {
	GetSchemaById(id int64) (*types.Schema, error)
	CreateSchema(schema *types.Schema) error
	UpdateSchema(schema *types.Schema) error
	DeleteSchema(id int64) error
}

type DBSchemaRepository struct {
	DB *sqlx.DB
}

func (repo DBSchemaRepository) GetSchemaById(id int64) (*types.Schema, error) {
	schema := types.Schema{}
	err := repo.DB.Get(&schema, "select * from schema where id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &schema, nil
	}
}

func (repo DBSchemaRepository) CreateSchema(schema *types.Schema) error {
	logrus.Trace("db.CreateSchema", schema)
	query := `
		insert into
		schema(
			created, user_id, name, description, complete,
			max_x, max_y, max_z, part_length,
			total_size, total_parts, license,
			search_tokens
		)
		values(
			:created, :user_id, :name, :description, :complete,
			:max_x, :max_y, :max_z, :part_length,
			:total_size, :total_parts, :license,
			to_tsvector(:description)
		)
		returning id
	`
	stmt, err := repo.DB.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&schema.ID, schema)
}

func (repo DBSchemaRepository) UpdateSchema(schema *types.Schema) error {
	query := `
		update schema
		set
			name = :name,
			description = :description,
			search_tokens = to_tsvector(:description),
			user_id = :user_id,
			license = :license,
			complete = :complete
		where id = :id
	`
	_, err := repo.DB.NamedExec(query, schema)
	return err
}

func (repo DBSchemaRepository) DeleteSchema(id int64) error {
	_, err := repo.DB.Exec("delete from schema where id = $1", id)
	return err
}
