package db

import (
	"blockexchange/types"
	"context"

	"github.com/vingarcia/ksql"
)

var schemaPartTable = ksql.NewTable("schemapart", "uid")

type SchemaPartRepository struct {
	kdb ksql.Provider
}

func (r *SchemaPartRepository) CreateOrUpdateSchemaPart(part *types.SchemaPart) error {
	return r.kdb.Transaction(context.Background(), func(p ksql.Provider) error {
		sp := &types.SchemaPart{}
		err := p.QueryOne(context.Background(), sp,
			"from schemapart where schema_uid = $1 and offset_x = $2 and offset_y = $3 and offset_z = $4",
			part.SchemaUID, part.OffsetX, part.OffsetY, part.OffsetZ,
		)
		if err == ksql.ErrRecordNotFound {
			// insert
			part.OrderID = int64(types.GetSchemaPartOrderID(part.OffsetX, part.OffsetY, part.OffsetZ))
			return p.Insert(context.Background(), schemaPartTable, part)
		} else if err == nil {
			// update
			sp.Mtime = part.Mtime
			sp.Data = part.Data
			sp.MetaData = part.MetaData
			return p.Patch(context.Background(), schemaPartTable, sp)
		} else {
			// err
			return err
		}
	})
}

func (r *SchemaPartRepository) GetBySchemaUIDAndOffset(schema_uid string, offset_x, offset_y, offset_z int) (*types.SchemaPart, error) {
	sp := &types.SchemaPart{}
	err := r.kdb.QueryOne(context.Background(), sp,
		"from schemapart where schema_uid = $1 and offset_x = $2 and offset_y = $3 and offset_z = $4",
		schema_uid, offset_x, offset_y, offset_z,
	)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	}
	return sp, err
}

func (r *SchemaPartRepository) GetBySchemaUIDAndRange(schema_uid string, x1, y1, z1, x2, y2, z2 int) ([]*types.SchemaPart, error) {
	list := []*types.SchemaPart{}
	q := `from schemapart
		where schema_uid = $1
		and offset_x >= $2
		and offset_y >= $3
		and offset_z >= $4
		and offset_x <= $5
		and offset_y <= $6
		and offset_z <= $7
	`
	return list, r.kdb.Query(context.Background(), &list, q, schema_uid, x1, y1, z1, x2, y2, z2)
}

func (r *SchemaPartRepository) RemoveBySchemaUIDAndOffset(schema_uid string, offset_x, offset_y, offset_z int) error {
	q := `
		delete from schemapart
		where schema_uid = $1
		and offset_x = $2
		and offset_y = $3
		and offset_z = $4
	`
	_, err := r.kdb.Exec(context.Background(), q, schema_uid, offset_x, offset_y, offset_z)
	return err
}

func (r *SchemaPartRepository) GetNextBySchemaUIDAndOffset(schema_uid string, offset_x, offset_y, offset_z int) (*types.SchemaPart, error) {
	order_id := types.GetSchemaPartOrderID(offset_x, offset_y, offset_z)
	sp := &types.SchemaPart{}
	q := `from schemapart
		where schema_uid = $1
		and order_id > $2
		order by order_id asc
		limit 1
	`
	err := r.kdb.QueryOne(context.Background(), sp, q, schema_uid, order_id)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	}
	return sp, err
}

func (r *SchemaPartRepository) GetNextBySchemaUIDAndMtime(schema_uid string, mtime int64) (*types.SchemaPart, error) {
	sp := &types.SchemaPart{}
	q := `from schemapart
		where schema_uid = $1
		and mtime > $2
		order by mtime asc limit 1
	`
	err := r.kdb.QueryOne(context.Background(), sp, q, schema_uid, mtime)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	}
	return sp, err
}

func (r *SchemaPartRepository) CountNextBySchemaUIDAndMtime(schema_uid string, mtime int64) (int64, error) {
	c := &types.Count{}
	return c.Count, r.kdb.QueryOne(context.Background(), c, "select count(*) as count from schemapart where schema_uid = $1 and mtime > $2", schema_uid, mtime)
}

func (r *SchemaPartRepository) GetFirstBySchemaUID(schema_uid string) (*types.SchemaPart, error) {
	sp := &types.SchemaPart{}
	q := `from schemapart
		where schema_uid = $1
		order by order_id asc
		limit 1
	`
	err := r.kdb.QueryOne(context.Background(), sp, q, schema_uid)
	if err == ksql.ErrRecordNotFound {
		return nil, nil
	}
	return sp, err
}
