package db

import (
	"blockexchange/types"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type SchemaRepository interface {
	GetSchemaById(id int64) (*types.Schema, error)
	GetSchemaByUsernameAndName(username, schemaname string) (*types.Schema, error)
	CreateSchema(schema *types.Schema) error
	UpdateSchema(schema *types.Schema) error
	DeleteSchema(id, user_id int64) error
	IncrementDownloads(id int64) error
	DeleteIncompleteSchema(user_id int64, name string) error
	DeleteOldIncompleteSchema(time_before int64) error
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

func (repo DBSchemaRepository) GetSchemaByUsernameAndName(username, schemaname string) (*types.Schema, error) {
	schema := types.Schema{}
	err := repo.DB.Get(&schema, "select * from schema where user_id = (select id from public.user where name = $1) and name = $2", username, schemaname)
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
			size_x, size_y, size_z,
			part_length,
			total_size, total_parts, license,
			search_tokens
		)
		values(
			:created, :user_id, :name, :description, :complete,
			:size_x, :size_y, :size_z,
			:part_length,
			:total_size, :total_parts, :license,
			to_tsvector(:description || ' ' || :name)
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
			search_tokens = to_tsvector(:description || ' ' || :name),
			mtime = :mtime,
			user_id = :user_id,
			license = :license,
			complete = :complete,
			downloads = :downloads
		where id = :id
	`
	_, err := repo.DB.NamedExec(query, schema)
	return err
}

func (repo DBSchemaRepository) IncrementDownloads(id int64) error {
	query := `
		update schema
		set downloads = downloads + 1
		where id = $1
	`
	_, err := repo.DB.Exec(query, id)
	return err
}

func (repo DBSchemaRepository) DeleteSchema(id, user_id int64) error {
	_, err := repo.DB.Exec("delete from schema where id = $1 and user_id = $2", id, user_id)
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

func (repo DBSchemaRepository) DeleteOldIncompleteSchema(time_before int64) error {
	q := `
		delete from schema where
			created < $1 and
			complete = false
	`
	_, err := repo.DB.Exec(q, time_before)
	return err
}

func (repo DBSchemaRepository) CalculateStats(id int64) error {
	q := `
		update schema s
		set total_size = (
			select
			coalesce(sum(length(data)) + sum(length(metadata)), 0)
			from schemapart sp where sp.schema_id = s.id
		),
		total_parts = (select count(*) from schemapart sp where sp.schema_id = s.id)
		where id = $1
	`
	_, err := repo.DB.Exec(q, id)
	return err
}
