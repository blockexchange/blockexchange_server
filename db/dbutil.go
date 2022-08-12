package db

import (
	"database/sql"
	"fmt"
	"strings"
)

type Entity interface {
	Columns() []string
	Table() string
}

type Selectable interface {
	Entity
	Scan(r func(dest ...any) error) error
}

type Insertable interface {
	Entity
	Values() []any
}

func Insert(d *sql.DB, entity Insertable, additionalStmts ...string) error {
	cols := entity.Columns()
	placeholders := make([]string, len(cols))
	for i := range cols {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	_, err := d.Exec(fmt.Sprintf(
		"insert into %s(%s) values(%s) %s",
		entity.Table(), strings.Join(cols, ","), strings.Join(placeholders, ","), strings.Join(additionalStmts, " ")),
		entity.Values()...,
	)

	return err
}

func Select[E Selectable](d *sql.DB, entity E, constraints string, params ...any) (E, error) {
	row := d.QueryRow(fmt.Sprintf(
		"select %s from %s %s",
		strings.Join(entity.Columns(), ","), entity.Table(), constraints),
		params...,
	)
	err := entity.Scan(row.Scan)
	return entity, err
}

func SelectMulti[E Selectable](d *sql.DB, p func() E, constraints string) ([]E, error) {
	entity := p()
	rows, err := d.Query(fmt.Sprintf("select %s from %s %s", strings.Join(entity.Columns(), ","), entity.Table(), constraints))
	if err != nil {
		return nil, err
	}

	list := make([]E, 0)
	for rows.Next() {
		entry := p()
		err = entry.Scan(rows.Scan)
		if err != nil {
			return nil, err
		}

		list = append(list, entry)
	}

	return list, nil
}
