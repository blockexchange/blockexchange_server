package db

import "database/sql"

type DBType interface {
	Columns() string
	Parameters() string
	Values() []any
	Scan(*sql.Row) error
}
