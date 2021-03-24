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
	DeleteIncompleteSchema(user_id int64, name string) error
	CalculateStats(id int64) error
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
			size_x, size_y, size_z, part_length,
			total_size, total_parts, license,
			search_tokens
		)
		values(
			:created, :user_id, :name, :description, :complete,
			:size_x, :size_y, :size_z, :part_length,
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
			complete = :complete,
			downloads = :downloads
		where id = :id
	`
	_, err := repo.DB.NamedExec(query, schema)
	return err
}

func (repo DBSchemaRepository) DeleteSchema(id int64) error {
	_, err := repo.DB.Exec("delete from schema where id = $1", id)
	return err
}

func (repo DBSchemaRepository) DeleteIncompleteSchema(user_id int64, name string) error {
	q := `
		delete from schema where
			user_id = $1 and
			name = $2 and
			complete = false
	`
	_, err := repo.DB.Exec(q, user_id, name)
	return err
}

func (repo DBSchemaRepository) CalculateStats(id int64) error {
	q := `
		update schema s
		set total_size = (select sum(length(data)) + sum(length(metadata)) from schemapart sp where sp.schema_id = s.id),
		total_parts = (select count(*) from schemapart sp where sp.schema_id = s.id)
		where id = $1
	`
	_, err := repo.DB.Exec(q, id)
	return err
}
