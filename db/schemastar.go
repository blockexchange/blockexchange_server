package db

import (
	"blockexchange/types"
	"context"

	"github.com/vingarcia/ksql"
)

var schemaStarTable = ksql.NewTable("user_schema_star", "user_uid", "schema_uid")

type SchemaStarRepository struct {
	kdb ksql.Provider
}

func (r *SchemaStarRepository) Create(ss *types.SchemaStar) error {
	return r.kdb.Insert(context.Background(), schemaStarTable, ss)
}

func (r *SchemaStarRepository) Delete(ss *types.SchemaStar) error {
	return r.kdb.Delete(context.Background(), schemaStarTable, ss)
}

func (r *SchemaStarRepository) GetBySchemaUID(schema_uid string) ([]*types.SchemaStar, error) {
	list := []*types.SchemaStar{}
	return list, r.kdb.Query(context.Background(), &list, "from user_schema_star where schema_uid = $1", schema_uid)
}

func (r *SchemaStarRepository) GetBySchemaAndUserID(schema_uid string, user_uid string) (*types.SchemaStar, error) {
	ss := &types.SchemaStar{}
	err := r.kdb.QueryOne(context.Background(), ss, "from user_schema_star where schema_uid = $1 and user_uid = $2", schema_uid, user_uid)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	} else {
		return ss, err
	}
}

func (r *SchemaStarRepository) CountBySchemaUID(schema_uid string) (int64, error) {
	c := &types.Count{}
	err := r.kdb.QueryOne(context.Background(), c, "select count(*) as count from user_schema_star where schema_uid = $1", schema_uid)
	return c.Count, err
}

func (r *SchemaStarRepository) CountByUserUID(user_uid string) (int64, error) {
	c := &types.Count{}
	err := r.kdb.QueryOne(context.Background(), c, "select count(*) from user_schema_star uss join schema s on s.uid = uss.schema_uid where s.user_uid = $1", user_uid)
	return c.Count, err
}
