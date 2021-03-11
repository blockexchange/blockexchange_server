package db

import "database/sql"

type Processor func()

type Repository struct {
	table string
}

func (r Repository) GetBy(db sql.DB, wherepart string, target interface{}) {

}

func NewRepository(table string) Repository {
	return Repository{
		table: table,
	}
}
