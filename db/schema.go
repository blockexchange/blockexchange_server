package db

import (
	"blockexchange/types"
	"database/sql"

	"github.com/minetest-go/dbutil"
)

type SchemaRepository struct {
	DB *sql.DB
}

func (repo SchemaRepository) GetSchemaById(id int64) (*types.Schema, error) {
	schema, err := dbutil.Select(repo.DB, &types.Schema{}, "where id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return schema, err
	}
}

func (repo SchemaRepository) GetSchemaByUserIDAndName(user_id int64, name string) (*types.Schema, error) {
	schema, err := dbutil.Select(repo.DB, &types.Schema{}, "where user_id = $1 and name = $2", user_id, name)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return schema, err
	}
}

func (repo SchemaRepository) GetSchemaByUsernameAndName(username, schemaname string) (*types.Schema, error) {
	schema, err := dbutil.Select(repo.DB, &types.Schema{}, "where user_id = (select id from public.user where name = $1) and name = $2", username, schemaname)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return schema, err
	}
}

func (repo SchemaRepository) CreateSchema(schema *types.Schema) error {
	err := dbutil.InsertReturning(repo.DB, schema, "id", &schema.ID)
	if err != nil {
		return err
	}
	return repo.UpdateSearchTokens(schema.ID)
}

func (repo SchemaRepository) UpdateSchema(schema *types.Schema) error {
	err := dbutil.Update(repo.DB, schema, "where id = $1", schema.ID)
	if err != nil {
		return err
	}
	return repo.UpdateSearchTokens(schema.ID)
}

func (repo SchemaRepository) UpdateSearchTokens(id int64) error {
	_, err := repo.DB.Exec(`
		update schema s
		set search_tokens = to_tsvector(s.description || ' ' || s.name)
		where id = $1
	`, id)
	return err
}

func (repo SchemaRepository) IncrementViews(id int64) error {
	query := `
		update schema
		set views = views + 1
		where id = $1
	`
	_, err := repo.DB.Exec(query, id)
	return err
}

func (repo SchemaRepository) IncrementDownloads(id int64) error {
	query := `
		update schema
		set downloads = downloads + 1
		where id = $1
	`
	_, err := repo.DB.Exec(query, id)
	return err
}

func (repo SchemaRepository) DeleteSchema(id int64) error {
	_, err := repo.DB.Exec("delete from schema where id = $1", id)
	return err
}

func (repo SchemaRepository) DeleteIncompleteSchema(user_id int64, name string) error {
	q := `
		delete from schema where
			user_id = $1 and
			name = $2 and
			complete = false
	`
	_, err := repo.DB.Exec(q, user_id, name)
	return err
}

func (repo SchemaRepository) DeleteOldIncompleteSchema(time_before int64) error {
	q := `
		delete from schema where
			created < $1 and
			complete = false
	`
	_, err := repo.DB.Exec(q, time_before)
	return err
}

func (repo SchemaRepository) CalculateStats(id int64) error {
	q := `
		update schema s
		set total_size = (
			select
			coalesce(sum(length(data)) + sum(length(metadata)), 0)
			from schemapart sp where sp.schema_id = s.id
		),
		total_parts = (select count(*) from schemapart sp where sp.schema_id = s.id),
		stars = (select count(*) from user_schema_star where schema_id = $1)
		where id = $1
	`
	_, err := repo.DB.Exec(q, id)
	return err
}
