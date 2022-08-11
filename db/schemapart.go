package db

import (
	"blockexchange/types"
	"database/sql"
	"fmt"
)

type SchemaPartRepository interface {
	CreateOrUpdateSchemaPart(part *types.SchemaPart) error
	GetBySchemaIDAndOffset(schema_id int64, offset_x, offset_y, offset_z int) (*types.SchemaPart, error)
	GetBySchemaIDAndRange(schema_id int64, x1, y1, z1, x2, y2, z2 int) ([]*types.SchemaPart, error)
	RemoveBySchemaIDAndOffset(schema_id int64, offset_x, offset_y, offset_z int) error
	GetNextBySchemaIDAndOffset(schema_id int64, offset_x, offset_y, offset_z int) (*types.SchemaPart, error)
	GetNextBySchemaIDAndMtime(schema_id int64, mtime int64) (*types.SchemaPart, error)
	GetFirstBySchemaID(schema_id int64) (*types.SchemaPart, error)
}

type DBSchemaPartRepository struct {
	DB *sql.DB
}

func (repo DBSchemaPartRepository) CreateOrUpdateSchemaPart(part *types.SchemaPart) error {
	conflict := `
		on conflict on constraint schemapart_unique_coords
		do update set data = EXCLUDED.data, metadata = EXCLUDED.metadata, mtime = EXCLUDED.mtime
	`
	_, err := repo.DB.Exec(
		fmt.Sprintf(`insert into schemapart(%s) values(%s) %s`, part.Columns(), part.Parameters(), conflict),
		part.Values()...,
	)
	return err
}

func (repo DBSchemaPartRepository) GetBySchemaIDAndOffset(schema_id int64, offset_x, offset_y, offset_z int) (*types.SchemaPart, error) {
	where := `
		where schema_id = $1
		and offset_x = $2
		and offset_y = $3
		and offset_z = $4
	`

	p := &types.SchemaPart{}
	row := repo.DB.QueryRow(
		fmt.Sprintf("select %s from schemapart %s", p.Columns(), where),
		schema_id, offset_x, offset_y, offset_z,
	)

	err := p.Scan(row.Scan)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (repo DBSchemaPartRepository) GetBySchemaIDAndRange(schema_id int64, x1, y1, z1, x2, y2, z2 int) ([]*types.SchemaPart, error) {
	list := []*types.SchemaPart{}
	sp := &types.SchemaPart{}
	where := `
		select *
		from schemapart
		where schema_id = $1
		and offset_x >= $2
		and offset_y >= $3
		and offset_z >= $4
		and offset_x <= $5
		and offset_y <= $6
		and offset_z <= $7
	`

	rows, err := repo.DB.Query(
		fmt.Sprintf("select %s from schemapart %s", sp.Columns(), where),
		schema_id, x1, y1, z1, x2, y2, z2,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		sp = &types.SchemaPart{}
		err = sp.Scan(rows.Scan)
		if err != nil {
			return nil, err
		}
		list = append(list, sp)
	}

	return list, nil
}

func (repo DBSchemaPartRepository) RemoveBySchemaIDAndOffset(schema_id int64, offset_x, offset_y, offset_z int) error {
	query := `
		delete
		from schemapart
		where schema_id = $1
		and offset_x = $2
		and offset_y = $3
		and offset_z = $4
	`
	_, err := repo.DB.Exec(query, schema_id, offset_x, offset_y, offset_z)
	return err
}

func (repo DBSchemaPartRepository) GetNextBySchemaIDAndOffset(schema_id int64, offset_x, offset_y, offset_z int) (*types.SchemaPart, error) {
	sp := &types.SchemaPart{}

	where := `
		where id > (
			select id from schemapart
			where schema_id = $1
			and offset_x = $2
			and offset_y = $3
			and offset_z = $4
		)
		and schema_id = $1
	`

	row := repo.DB.QueryRow(
		fmt.Sprintf("select %s from schemapart %s order by id asc limit 1", sp.Columns(), where),
		schema_id, offset_x, offset_y, offset_z,
	)

	err := sp.Scan(row.Scan)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return sp, nil
}

func (repo DBSchemaPartRepository) GetNextBySchemaIDAndMtime(schema_id int64, mtime int64) (*types.SchemaPart, error) {
	sp := &types.SchemaPart{}

	where := `
		where mtime > $2
		and schema_id = $1
	`

	row := repo.DB.QueryRow(
		fmt.Sprintf("select %s from schemapart %s order by mtime asc limit 1", sp.Columns(), where),
		schema_id, mtime,
	)

	err := sp.Scan(row.Scan)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return sp, nil
}

func (repo DBSchemaPartRepository) GetFirstBySchemaID(schema_id int64) (*types.SchemaPart, error) {
	sp := &types.SchemaPart{}

	where := `
		where id = (
			select min(id) from schemapart
			where schema_id = $1
		)
		and schema_id = $1
	`

	row := repo.DB.QueryRow(
		fmt.Sprintf("select %s from schemapart %s limit 1", sp.Columns(), where),
		schema_id,
	)

	err := sp.Scan(row.Scan)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return sp, nil
}
