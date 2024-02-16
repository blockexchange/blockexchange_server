package db

import (
	"blockexchange/types"
	"context"

	"github.com/google/uuid"
	"github.com/vingarcia/ksql"
)

var schemaTable = ksql.NewTable("schema", "uid")

type SchemaRepository struct {
	kdb ksql.Provider
}

func (r *SchemaRepository) GetSchemaByUID(uid string) (*types.Schema, error) {
	s := &types.Schema{}
	err := r.kdb.QueryOne(context.Background(), s, "from schema where uid = $1", uid)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	} else {
		return s, err
	}
}

func (r *SchemaRepository) GetSchemaByUserUIDAndName(user_uid string, name string) (*types.Schema, error) {
	s := &types.Schema{}
	err := r.kdb.QueryOne(context.Background(), s, "from schema where user_uid = $1 and name = $2", user_uid, name)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	} else {
		return s, err
	}
}

func (r *SchemaRepository) GetSchemaByUsernameAndName(username, schemaname string) (*types.Schema, error) {
	s := &types.Schema{}
	err := r.kdb.QueryOne(context.Background(), s, "from schema where user_uid = (select uid from public.user where name = $1) and name = $2", username, schemaname)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	} else {
		return s, err
	}
}

func (r *SchemaRepository) CreateSchema(schema *types.Schema) error {
	if schema.UID == "" {
		schema.UID = uuid.NewString()
	}
	return r.kdb.Insert(context.Background(), schemaTable, schema)
}

func (r *SchemaRepository) UpdateSchema(schema *types.Schema) error {
	return r.kdb.Patch(context.Background(), schemaTable, schema)
}

func (r *SchemaRepository) IncrementViews(uid string) error {
	query := `
		update schema
		set views = views + 1
		where uid = $1
	`
	_, err := r.kdb.Exec(context.Background(), query, uid)
	return err
}

func (r *SchemaRepository) IncrementDownloads(uid string) error {
	query := `
		update schema
		set downloads = downloads + 1
		where uid = $1
	`
	_, err := r.kdb.Exec(context.Background(), query, uid)
	return err
}

func (r *SchemaRepository) DeleteSchema(uid string) error {
	return r.kdb.Delete(context.Background(), schemaTable, uid)
}

func (r *SchemaRepository) DeleteIncompleteSchema(user_uid string, name string) error {
	q := `
		delete from schema where
			user_uid = $1 and
			name = $2 and
			complete = false
	`
	_, err := r.kdb.Exec(context.Background(), q, user_uid, name)
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

func (r *SchemaRepository) CalculateStats(uid string) error {
	q := `
		update schema s
		set total_size = (
			select
			coalesce(sum(length(data)) + sum(length(metadata)), 0)
			from schemapart sp where sp.schema_uid = s.uid
		),
		total_parts = (select count(*) from schemapart sp where sp.schema_uid = s.uid),
		stars = (select count(*) from user_schema_star where schema_uid = $1)
		where uid = $1
	`
	_, err := r.kdb.Exec(context.Background(), q, uid)
	return err
}

func (r *SchemaRepository) GetTotalSize() (int64, error) {
	c := &types.Count{}
	err := r.kdb.QueryOne(context.Background(), c, "select sum(total_size) as count from schema")
	return c.Count, err
}
