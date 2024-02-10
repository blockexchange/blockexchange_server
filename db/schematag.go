package db

import (
	"blockexchange/types"
	"context"

	"github.com/google/uuid"
	"github.com/vingarcia/ksql"
)

var schemaTagTable = ksql.NewTable("schematag", "uid")

type SchemaTagRepository struct {
	kdb ksql.Provider
}

func (r *SchemaTagRepository) Create(st *types.SchemaTag) error {
	if st.UID == "" {
		st.UID = uuid.NewString()
	}
	return r.kdb.Insert(context.Background(), schemaTagTable, st)
}

func (r *SchemaTagRepository) Delete(schema_id int64, tag_uid string) error {
	_, err := r.kdb.Exec(context.Background(), "delete from schematag where schema_id = $1 and tag_uid = $2", schema_id, tag_uid)
	return err
}

func (r *SchemaTagRepository) GetBySchemaID(schema_id int64) ([]*types.SchemaTag, error) {
	list := []*types.SchemaTag{}
	err := r.kdb.Query(context.Background(), &list, "from schematag where schema_id = $1", schema_id)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	} else {
		return list, err
	}
}

func (r *SchemaTagRepository) GetBySchemaIDs(schema_ids []int64) ([]*types.SchemaTag, error) {
	list := []*types.SchemaTag{}
	err := r.kdb.Query(context.Background(), &list, "from schematag where schema_id = any($1::bigint[])", schema_ids)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	} else {
		return list, err
	}
}

func (r *SchemaTagRepository) GetByTagUID(tag_uid int64) ([]*types.SchemaTag, error) {
	list := []*types.SchemaTag{}
	err := r.kdb.Query(context.Background(), &list, "from schematag where tag_uid = $1", tag_uid)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	} else {
		return list, err
	}
}
