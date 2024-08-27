package core

import (
	"blockexchange/db"
	"blockexchange/types"

	mt "github.com/minetest-go/types"
)

type SchemaPartMap struct {
	repo          *db.SchemaPartRepository
	schema        *types.Schema
	blocks        map[int64]*types.SchemaPart
	changedblocks map[int64]bool
}

func NewSchemaPartMap(repo *db.SchemaPartRepository, schema *types.Schema) *SchemaPartMap {
	return &SchemaPartMap{
		repo:          repo,
		schema:        schema,
		blocks:        make(map[int64]*types.SchemaPart),
		changedblocks: make(map[int64]bool),
	}
}

func (m *SchemaPartMap) SetNode(pos *mt.Pos, node *mt.Node) error {
	// TODO
	return nil
}

func (m *SchemaPartMap) GetNode(pos *mt.Pos) (*mt.Node, error) {
	// TODO
	return nil, nil
}

func (m *SchemaPartMap) Close() error {
	// TODO
	return nil
}
