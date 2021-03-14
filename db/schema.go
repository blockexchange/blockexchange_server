package db

import (
	"blockexchange/types"

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
	list := []types.Schema{}
	err := repo.DB.Select(&list, "select * from schema where id = $1", id)
	if err != nil {
		return nil, err
	} else if len(list) == 1 {
		return &list[0], nil
	} else {
		return nil, nil
	}
}

func (repo DBSchemaRepository) CreateSchema(schema *types.Schema) error {
	logrus.Trace("db.CreateSchema", schema)
	query := `
		insert into
		schema(
			created, user_id, name, description, complete,
			max_x, max_y, max_z, part_length,
			total_size, total_parts, downloads, license
		)
		values(
			:created, :user_id, :name, :description, :complete,
			:max_x, :max_y, :max_z, :part_length,
			:total_size, :total_parts, :downloads, :license
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
			search_tokens = to_tsvector(concat(:name, ' ', :description)),
			user_id = :user_id,
			license = :license
		where id = :id
	`
	_, err := repo.DB.NamedExec(query, schema)
	return err
}

func (repo DBSchemaRepository) DeleteSchema(id int64) error {
	_, err := repo.DB.Exec("delete from schema where id = $1", id)
	return err
}
