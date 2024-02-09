package db

import (
	"blockexchange/types"
	"context"

	"github.com/vingarcia/ksql"
)

type MetaRepository struct {
	kdb ksql.Provider
}

func (r *MetaRepository) CountEntries(table string) (int64, error) {
	c := &types.Count{}
	err := r.kdb.QueryOne(context.Background(), c, "SELECT reltuples::bigint as count FROM pg_catalog.pg_class WHERE relname = $1", table)
	return c.Count, err
}
