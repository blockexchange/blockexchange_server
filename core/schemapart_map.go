package core

import (
	"blockexchange/db"
	"blockexchange/parser"
	"blockexchange/types"
	"fmt"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	mt "github.com/minetest-go/types"
)

type SchemaPartMap struct {
	repo         *db.SchemaPartRepository
	schema       *types.Schema
	cache        *expirable.LRU[int64, *parser.ParsedSchemaPart]
	changedparts map[int64]*parser.ParsedSchemaPart
}

func NewSchemaPartMap(repo *db.SchemaPartRepository, schema *types.Schema) *SchemaPartMap {
	return &SchemaPartMap{
		repo:         repo,
		schema:       schema,
		cache:        expirable.NewLRU[int64, *parser.ParsedSchemaPart](1000, nil, time.Second*10),
		changedparts: make(map[int64]*parser.ParsedSchemaPart),
	}
}

func (m *SchemaPartMap) getKey(mbpos *mt.Pos) int64 {
	return int64(mbpos[0]) +
		int64(mbpos[1])<<16 +
		int64(mbpos[2])<<32
}

func (m *SchemaPartMap) getPart(mbpos *mt.Pos) (*parser.ParsedSchemaPart, error) {
	key := m.getKey(mbpos)

	if m.changedparts[key] != nil {
		// found block in changed map
		return m.changedparts[key], nil
	}

	mapblock, found := m.cache.Get(key)
	if found && mapblock == airOnlyMapblock {
		return nil, nil
	}

	if !found {
		schemapart, err := m.repo.GetBySchemaUIDAndOffset(m.schema.UID, mbpos[0], mbpos[1], mbpos[2])
		if err != nil {
			return nil, fmt.Errorf("get schemapart error @ %s: %v", mbpos, err)
		}

		if schemapart != nil {
			mapblock, err = parser.ParseSchemaPart(schemapart)
			if err != nil {
				return nil, fmt.Errorf("parse error @ %s: %v", mbpos, err)
			}

			if len(mapblock.Meta.NodeMapping) == 1 && mapblock.Meta.NodeMapping["air"] > 0 {
				// mark as air-only
				m.cache.Add(key, airOnlyMapblock)
				return nil, nil
			}

			m.cache.Add(key, mapblock)
		} else {
			// not found, mark as air-only
			m.cache.Add(key, airOnlyMapblock)
			return nil, nil
		}
	}

	return mapblock, nil
}

func (m *SchemaPartMap) SetNode(pos *mt.Pos, node *mt.Node) error {
	mbpos := pos.Divide(16).Multiply(16)

	mapblock, err := m.getPart(mbpos)
	if err != nil {
		return fmt.Errorf("getpart error: %v", err)
	}
	if mapblock == nil {
		// TODO: create new mapblock
		return nil
	}

	rel_pos := pos.Subtract(mbpos)
	index := mapblock.GetIndex(rel_pos.X(), rel_pos.Y(), rel_pos.Z())
	if index >= len(mapblock.NodeIDS) {
		return fmt.Errorf("index mismatch: got %d, length: %d, rel_pos: %s, abs_pos: %s", index, len(mapblock.NodeIDS), rel_pos, pos)
	}
	nodeid := mapblock.NodeIDS[index]

	if mapblock.Meta.NodeMapping[node.Name] == 0 {
		// add to mapping
		maxnodeid := int16(0)
		for nodeid := range mapblock.NodeNameLookup {
			maxnodeid = max(nodeid, maxnodeid)
		}
		nodeid = maxnodeid + 1
		mapblock.NodeNameLookup[nodeid] = node.Name
		mapblock.Meta.NodeMapping[node.Name] = nodeid
	}

	mapblock.Param1[index] = byte(node.Param1)
	mapblock.Param2[index] = byte(node.Param2)
	mapblock.NodeIDS[index] = nodeid

	key := m.getKey(mbpos)
	m.changedparts[key] = mapblock

	// TODO: meta and stuff
	return nil
}

func (m *SchemaPartMap) GetNode(pos *mt.Pos) (*mt.Node, error) {
	mbpos := pos.Divide(16).Multiply(16)

	mapblock, err := m.getPart(mbpos)
	if err != nil {
		return nil, fmt.Errorf("getpart error: %v", err)
	}
	if mapblock == nil {
		return nil, nil
	}

	rel_pos := pos.Subtract(mbpos)
	index := mapblock.GetIndex(rel_pos.X(), rel_pos.Y(), rel_pos.Z())
	if index >= len(mapblock.NodeIDS) {
		return nil, fmt.Errorf("index mismatch: got %d, length: %d, rel_pos: %s, abs_pos: %s", index, len(mapblock.NodeIDS), rel_pos, pos)
	}
	nodeid := mapblock.NodeIDS[index]

	return &mt.Node{
		Name:   mapblock.NodeNameLookup[nodeid],
		Param1: int(mapblock.Param1[index]),
		Param2: int(mapblock.Param2[index]),
		Pos:    pos,
	}, nil
}

func (m *SchemaPartMap) Close() error {
	// TODO: save all changed parts
	return nil
}
