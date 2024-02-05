package db

import "database/sql"

type MetaRepository struct {
	db *sql.DB
}

func (r *MetaRepository) CountEntries(table string) (int64, error) {
	row := r.db.QueryRow("SELECT reltuples::bigint FROM pg_catalog.pg_class WHERE relname = $1", table)
	var count int64
	return count, row.Scan(&count)
}
