package db

import (
	"blockexchange/types"
	"database/sql"
)

type SchemaPartRepository struct {
	DB *sql.DB
}

func (repo SchemaPartRepository) CreateOrUpdateSchemaPart(part *types.SchemaPart) error {
	conflict := `
		on conflict on constraint schemapart_unique_coords
		do update set data = EXCLUDED.data, metadata = EXCLUDED.metadata, mtime = EXCLUDED.mtime
	`
	return Insert(repo.DB, part, conflict)
}

func (repo SchemaPartRepository) GetBySchemaIDAndOffset(schema_id int64, offset_x, offset_y, offset_z int) (*types.SchemaPart, error) {
	where := `
		where schema_id = $1
		and offset_x = $2
		and offset_y = $3
		and offset_z = $4
	`

	sp, err := Select(repo.DB, &types.SchemaPart{}, where, schema_id, offset_x, offset_y, offset_z)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return sp, err
	}
}

func (repo SchemaPartRepository) GetBySchemaIDAndRange(schema_id int64, x1, y1, z1, x2, y2, z2 int) ([]*types.SchemaPart, error) {
	constraints := `
		where schema_id = $1
		and offset_x >= $2
		and offset_y >= $3
		and offset_z >= $4
		and offset_x <= $5
		and offset_y <= $6
		and offset_z <= $7
	`

	return SelectMulti(repo.DB, func() *types.SchemaPart { return &types.SchemaPart{} }, constraints)
}

func (repo SchemaPartRepository) RemoveBySchemaIDAndOffset(schema_id int64, offset_x, offset_y, offset_z int) error {
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

func (repo SchemaPartRepository) GetNextBySchemaIDAndOffset(schema_id int64, offset_x, offset_y, offset_z int) (*types.SchemaPart, error) {
	constraints := `
		where id > (
			select id from schemapart
			where schema_id = $1
			and offset_x = $2
			and offset_y = $3
			and offset_z = $4
		)
		and schema_id = $1
		order by id asc limit 1
	`
	sp, err := Select(repo.DB, &types.SchemaPart{}, constraints, schema_id, offset_x, offset_y, offset_z)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return sp, err
	}
}

func (repo SchemaPartRepository) GetNextBySchemaIDAndMtime(schema_id int64, mtime int64) (*types.SchemaPart, error) {
	constraints := `
		where mtime > $2
		and schema_id = $1
		order by mtime asc limit 1
	`
	sp, err := Select(repo.DB, &types.SchemaPart{}, constraints, schema_id, mtime)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return sp, err
	}
}

func (repo SchemaPartRepository) GetFirstBySchemaID(schema_id int64) (*types.SchemaPart, error) {
	constraints := `
		where id = (
			select min(id) from schemapart
			where schema_id = $1
		)
		and schema_id = $1
		limit 1
	`
	sp, err := Select(repo.DB, &types.SchemaPart{}, constraints, schema_id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return sp, err
	}
}
