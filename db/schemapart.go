package db

import (
	"blockexchange/types"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
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
	DB *sqlx.DB
}

func (repo DBSchemaPartRepository) CreateOrUpdateSchemaPart(part *types.SchemaPart) error {
	logrus.WithFields(logrus.Fields{
		"schema_id": part.SchemaID,
	}).Trace("db.CreateOrUpdateSchemaPart")

	query := `
		insert into
		schemapart(schema_id, offset_x, offset_y, offset_z, mtime, data, metadata)
		values(:schema_id, :offset_x, :offset_y, :offset_z, :mtime, :data, :metadata)
		on conflict on constraint schemapart_unique_coords
		do update set data = EXCLUDED.data, metadata = EXCLUDED.metadata, mtime = EXCLUDED.mtime
		returning id
	`
	stmt, err := repo.DB.PrepareNamed(query)
	if err != nil {
		return err
	}
	return stmt.Get(&part.ID, part)
}

func (repo DBSchemaPartRepository) GetBySchemaIDAndOffset(schema_id int64, offset_x, offset_y, offset_z int) (*types.SchemaPart, error) {
	list := []types.SchemaPart{}
	query := `
		select *
		from schemapart
		where schema_id = $1
		and offset_x = $2
		and offset_y = $3
		and offset_z = $4
	`
	err := repo.DB.Select(&list, query, schema_id, offset_x, offset_y, offset_z)
	if err != nil {
		return nil, err
	} else if len(list) == 1 {
		return &list[0], nil
	} else {
		return nil, nil
	}
}

func (repo DBSchemaPartRepository) GetBySchemaIDAndRange(schema_id int64, x1, y1, z1, x2, y2, z2 int) ([]*types.SchemaPart, error) {
	list := []*types.SchemaPart{}
	query := `
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
	err := repo.DB.Select(&list, query, schema_id, x1, y1, z1, x2, y2, z2)
	return list, err
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
	list := []types.SchemaPart{}
	query := `
		select *
		from schemapart
		where id > (
			select id from schemapart
			where schema_id = $1
			and offset_x = $2
			and offset_y = $3
			and offset_z = $4
		)
		and schema_id = $1
		order by id asc
		limit 1
	`
	err := repo.DB.Select(&list, query, schema_id, offset_x, offset_y, offset_z)
	if err != nil {
		return nil, err
	} else if len(list) == 1 {
		return &list[0], nil
	} else {
		return nil, nil
	}
}

func (repo DBSchemaPartRepository) GetNextBySchemaIDAndMtime(schema_id int64, mtime int64) (*types.SchemaPart, error) {
	list := []types.SchemaPart{}
	query := `
		select *
		from schemapart
		where mtime > $2
		and schema_id = $1
		order by mtime asc
		limit 1
	`
	err := repo.DB.Select(&list, query, schema_id, mtime)
	if err != nil {
		return nil, err
	} else if len(list) == 1 {
		return &list[0], nil
	} else {
		return nil, nil
	}
}

func (repo DBSchemaPartRepository) GetFirstBySchemaID(schema_id int64) (*types.SchemaPart, error) {
	list := []types.SchemaPart{}
	query := `
		select *
		from schemapart
		where id = (
			select min(id) from schemapart
			where schema_id = $1
		)
		and schema_id = $1
		limit 1
	`
	err := repo.DB.Select(&list, query, schema_id)
	if err != nil {
		return nil, err
	} else if len(list) == 1 {
		return &list[0], nil
	} else {
		return nil, nil
	}
}
