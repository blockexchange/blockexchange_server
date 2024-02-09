package db

import (
	"blockexchange/types"
	"context"

	"github.com/vingarcia/ksql"
)

var schemaTable = ksql.NewTable("schema", "id")

type SchemaRepository struct {
	kdb ksql.Provider
}

func (r *SchemaRepository) GetSchemaById(id int64) (*types.Schema, error) {
	s := &types.Schema{}
	err := r.kdb.QueryOne(context.Background(), s, "from schema where id = $1", id)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	} else {
		return s, err
	}
}

func (r *SchemaRepository) GetSchemaByUserIDAndName(user_id int64, name string) (*types.Schema, error) {
	s := &types.Schema{}
	err := r.kdb.QueryOne(context.Background(), s, "from schema where user_id = $1 and name = $2", user_id, name)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	} else {
		return s, err
	}
}

func (r *SchemaRepository) GetSchemaByUsernameAndName(username, schemaname string) (*types.Schema, error) {
	s := &types.Schema{}
	err := r.kdb.QueryOne(context.Background(), s, "from schema where user_id = (select id from public.user where name = $1) and name = $2", username, schemaname)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	} else {
		return s, err
	}
}

func (r *SchemaRepository) CreateSchema(schema *types.Schema) error {
	return r.kdb.Insert(context.Background(), schemaTable, schema)
}

func (r *SchemaRepository) UpdateSchema(schema *types.Schema) error {
	return r.kdb.Patch(context.Background(), schemaTable, schema)
}

func (r *SchemaRepository) IncrementViews(id int64) error {
	query := `
		update schema
		set views = views + 1
		where id = $1
	`
	_, err := r.kdb.Exec(context.Background(), query, id)
	return err
}

func (r *SchemaRepository) IncrementDownloads(id int64) error {
	query := `
		update schema
		set downloads = downloads + 1
		where id = $1
	`
	_, err := r.kdb.Exec(context.Background(), query, id)
	return err
}

func (r *SchemaRepository) DeleteSchema(id int64) error {
	return r.kdb.Delete(context.Background(), schemaTable, id)
}

func (r *SchemaRepository) DeleteIncompleteSchema(user_id int64, name string) error {
	q := `
		delete from schema where
			user_id = $1 and
			name = $2 and
			complete = false
	`
	_, err := r.kdb.Exec(context.Background(), q, user_id, name)
	return err
}

func (r *SchemaRepository) DeleteOldIncompleteSchema(time_before int64) error {
	q := `
		delete from schema where
			created < $1 and
			complete = false
	`
	_, err := r.kdb.Exec(context.Background(), q, time_before)
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
	_, err := r.kdb.Exec(context.Background(), q, id)
	return err
}

func (r *SchemaRepository) GetTotalSize() (int64, error) {
	c := &types.Count{}
	err := r.kdb.QueryOne(context.Background(), c, "select sum(total_size) as count from schema")
	return c.Count, err
}
