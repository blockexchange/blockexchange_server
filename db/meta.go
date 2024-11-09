package db

import (
	"gorm.io/gorm"
)

type MetaRepository struct {
	g *gorm.DB
}

func (r *MetaRepository) CountEntries(table string) (int64, error) {
	var c int64
	return c, r.g.Raw("SELECT reltuples::bigint as count FROM pg_catalog.pg_class WHERE relname = ?", table).Scan(&c).Error
}
