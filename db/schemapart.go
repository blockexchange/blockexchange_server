package db

import (
	"blockexchange/types"
	"database/sql"

	"github.com/minetest-go/dbutil"
)

func NewSchemaPartRepository(DB *sql.DB) *SchemaPartRepository {
	return &SchemaPartRepository{
		DB:  DB,
		dbu: dbutil.New(DB, dbutil.DialectPostgres, func() *types.SchemaPart { return &types.SchemaPart{} }),
	}
}

type SchemaPartRepository struct {
	DB  *sql.DB
	dbu *dbutil.DBUtil[*types.SchemaPart]
}

func (r *SchemaPartRepository) CreateOrUpdateSchemaPart(part *types.SchemaPart) error {
	conflict := `
		on conflict on constraint schemapart_unique_coords
		do update set data = EXCLUDED.data, metadata = EXCLUDED.metadata, mtime = EXCLUDED.mtime
	`
	return r.dbu.Insert(part, conflict)
}

func (r *SchemaPartRepository) GetBySchemaIDAndOffset(schema_id int64, offset_x, offset_y, offset_z int) (*types.SchemaPart, error) {
	where := `
		where schema_id = %s
		and offset_x = %s
		and offset_y = %s
		and offset_z = %s
	`
	sp, err := r.dbu.Select(where, schema_id, offset_x, offset_y, offset_z)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return sp, err
	}
}

func (r *SchemaPartRepository) GetBySchemaIDAndRange(schema_id int64, x1, y1, z1, x2, y2, z2 int) ([]*types.SchemaPart, error) {
	constraints := `
		where schema_id = %s
		and offset_x >= %s
		and offset_y >= %s
		and offset_z >= %s
		and offset_x <= %s
		and offset_y <= %s
		and offset_z <= %s
	`
	return r.dbu.SelectMulti(constraints, schema_id, x1, y1, z1, x2, y2, z2)
}

func (r *SchemaPartRepository) RemoveBySchemaIDAndOffset(schema_id int64, offset_x, offset_y, offset_z int) error {
	constraints := `
		where schema_id = %s
		and offset_x = %s
		and offset_y = %s
		and offset_z = %s
	`
	return r.dbu.Delete(constraints, schema_id, offset_x, offset_y, offset_z)
}

func (r *SchemaPartRepository) GetNextBySchemaIDAndOffset(schema_id int64, offset_x, offset_y, offset_z int) (*types.SchemaPart, error) {
	constraints := `
		where id > (
			select id from schemapart
			where schema_id = %s
			and offset_x = %s
			and offset_y = %s
			and offset_z = %s
		)
		and schema_id = %s
		order by id asc limit 1
	`
	sp, err := r.dbu.Select(constraints, schema_id, offset_x, offset_y, offset_z, schema_id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return sp, err
	}
}

func (r *SchemaPartRepository) GetNextBySchemaIDAndMtime(schema_id int64, mtime int64) (*types.SchemaPart, error) {
	constraints := `
		where schema_id = %s
		and mtime > %s
		order by mtime asc limit 1
	`
	sp, err := r.dbu.Select(constraints, schema_id, mtime)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return sp, err
	}
}

func (r *SchemaPartRepository) CountNextBySchemaIDAndMtime(schema_id int64, mtime int64) (int, error) {
	constraints := `
		where schema_id = %s
		and mtime > %s
	`
	return r.dbu.Count(constraints, schema_id, mtime)
}

func (r *SchemaPartRepository) GetFirstBySchemaID(schema_id int64) (*types.SchemaPart, error) {
	constraints := `
		where id = (
			select min(id) from schemapart
			where schema_id = %s
		)
		and schema_id = %s
		limit 1
	`
	sp, err := r.dbu.Select(constraints, schema_id, schema_id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return sp, err
	}
}
