package db

import (
	"blockexchange/types"
	"context"

	"github.com/vingarcia/ksql"
)

var schemaStarTable = ksql.NewTable("user_schema_star", "user_id", "schema_id")

type SchemaStarRepository struct {
	kdb ksql.Provider
}

func (r *SchemaStarRepository) Create(ss *types.SchemaStar) error {
	return r.kdb.Insert(context.Background(), schemaStarTable, ss)
}

func (r *SchemaStarRepository) Delete(ss *types.SchemaStar) error {
	return r.kdb.Delete(context.Background(), schemaStarTable, ss)
}

func (r *SchemaStarRepository) GetBySchemaID(schema_id int64) ([]*types.SchemaStar, error) {
	list := []*types.SchemaStar{}
	return list, r.kdb.Query(context.Background(), &list, "from user_schema_star where schema_id = $1", schema_id)
}

func (r *SchemaStarRepository) GetBySchemaAndUserID(schema_id int64, user_id int64) (*types.SchemaStar, error) {
	ss := &types.SchemaStar{}
	err := r.kdb.QueryOne(context.Background(), ss, "from user_schema_star where schema_id = $1 and user_id = $2", schema_id, user_id)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	} else {
		return ss, err
	}
}

func (r *SchemaStarRepository) CountBySchemaID(schema_id int64) (int64, error) {
	c := &types.Count{}
	err := r.kdb.QueryOne(context.Background(), c, "select count(*) as count from user_schema_star where schema_id = $1", schema_id)
	return c.Count, err
}

func (r *SchemaStarRepository) CountByUserID(user_id int64) (int64, error) {
	c := &types.Count{}
	err := r.kdb.QueryOne(context.Background(), c, "select count(*) from user_schema_star uss join schema s on s.id = uss.schema_id where s.user_id = $1", user_id)
	return c.Count, err
}
