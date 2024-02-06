package db

import (
	"blockexchange/types"
	"database/sql"

	"github.com/minetest-go/dbutil"
)

func NewSchemaRepository(DB *sql.DB) *SchemaRepository {
	return &SchemaRepository{
		DB:  DB,
		dbu: dbutil.New(DB, dbutil.DialectPostgres, func() *types.Schema { return &types.Schema{} }),
	}
}

type SchemaRepository struct {
	DB  *sql.DB
	dbu *dbutil.DBUtil[*types.Schema]
}

func (r *SchemaRepository) GetSchemaById(id int64) (*types.Schema, error) {
	schema, err := r.dbu.Select("where id = %s", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return schema, err
	}
}

func (r *SchemaRepository) GetSchemaByUserIDAndName(user_id int64, name string) (*types.Schema, error) {
	schema, err := r.dbu.Select("where user_id = %s and name = %s", user_id, name)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return schema, err
	}
}

func (r *SchemaRepository) GetSchemaByUsernameAndName(username, schemaname string) (*types.Schema, error) {
	schema, err := r.dbu.Select("where user_id = (select id from public.user where name = %s) and name = %s", username, schemaname)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return schema, err
	}
}

func (r *SchemaRepository) CreateSchema(schema *types.Schema) error {
	err := r.dbu.InsertReturning(schema, "id", &schema.ID)
	if err != nil {
		return err
	}
	return r.UpdateSearchTokens(schema.ID)
}

func (r *SchemaRepository) UpdateSchema(schema *types.Schema) error {
	err := r.dbu.Update(schema, "where id = %s", schema.ID)
	if err != nil {
		return err
	}
	return r.UpdateSearchTokens(schema.ID)
}

func (r *SchemaRepository) UpdateSearchTokens(id int64) error {
	_, err := r.DB.Exec(`
		update schema s
		set search_tokens = to_tsvector(s.description || ' ' || s.name)
		where id = $1
	`, id)
	return err
}

func (r *SchemaRepository) IncrementViews(id int64) error {
	query := `
		update schema
		set views = views + 1
		where id = $1
	`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *SchemaRepository) IncrementDownloads(id int64) error {
	query := `
		update schema
		set downloads = downloads + 1
		where id = $1
	`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *SchemaRepository) DeleteSchema(id int64) error {
	_, err := r.DB.Exec("delete from schema where id = $1", id)
	return err
}

func (r *SchemaRepository) DeleteIncompleteSchema(user_id int64, name string) error {
	q := `
		delete from schema where
			user_id = $1 and
			name = $2 and
			complete = false
	`
	_, err := r.DB.Exec(q, user_id, name)
	return err
}

func (r *SchemaRepository) DeleteOldIncompleteSchema(time_before int64) error {
	q := `
		delete from schema where
			created < $1 and
			complete = false
	`
	_, err := r.DB.Exec(q, time_before)
	return err
}

func (r *SchemaRepository) CalculateStats(id int64) error {
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
	_, err := r.DB.Exec(q, id)
	return err
}

func (r *SchemaRepository) GetTotalSize() (int64, error) {
	row := r.DB.QueryRow("select sum(total_size) from schema")
	var size int64
	return size, row.Scan(&size)
}
