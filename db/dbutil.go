package db

import (
	"database/sql"
	"fmt"
	"strings"
)

const (
	InsertAction = "insert"
	UpdateAction = "update"
	SelectAction = "select"
)

type Entity interface {
	Columns(action string) []string
	Table() string
}

type Selectable interface {
	Entity
	Scan(action string, r func(dest ...any) error) error
}

type Insertable interface {
	Entity
	Values(action string) []any
}

func Insert(d *sql.DB, entity Insertable, additionalStmts ...string) error {
	cols := entity.Columns(InsertAction)
	placeholders := make([]string, len(cols))
	for i := range cols {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	_, err := d.Exec(fmt.Sprintf(
		"insert into %s(%s) values(%s) %s",
		entity.Table(), strings.Join(cols, ","), strings.Join(placeholders, ","), strings.Join(additionalStmts, " ")),
		entity.Values(InsertAction)...,
	)

	return err
}

func InsertReturning(d *sql.DB, entity Insertable, retField string, retValue any) error {
	cols := entity.Columns(InsertAction)
	placeholders := make([]string, len(cols))
	for i := range cols {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	stmt, err := d.Prepare(fmt.Sprintf(
		"insert into %s(%s) values(%s) returning %s",
		entity.Table(), strings.Join(cols, ","), strings.Join(placeholders, ","), retField),
	)
	if err != nil {
		return err
	}

	row := stmt.QueryRow(entity.Values(InsertAction)...)
	err = row.Scan(retValue)

	return err
}

func Select[E Selectable](d *sql.DB, entity E, constraints string, params ...any) (E, error) {
	row := d.QueryRow(fmt.Sprintf(
		"select %s from %s %s",
		strings.Join(entity.Columns(SelectAction), ","), entity.Table(), constraints),
		params...,
	)
	err := entity.Scan(SelectAction, row.Scan)
	return entity, err
}

func SelectMulti[E Selectable](d *sql.DB, p func() E, constraints string, params ...any) ([]E, error) {
	entity := p()
	rows, err := d.Query(fmt.Sprintf("select %s from %s %s", strings.Join(entity.Columns(SelectAction), ","), entity.Table(), constraints), params...)
	if err != nil {
		return nil, err
	}

	list := make([]E, 0)
	for rows.Next() {
		entry := p()
		err = entry.Scan(SelectAction, rows.Scan)
		if err != nil {
			return nil, err
		}

		list = append(list, entry)
	}

	return list, nil
}
