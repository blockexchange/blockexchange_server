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

const (
	CONTENT_UNKNOWN = 125
	CONTENT_AIR     = 126
	CONTENT_IGNORE  = 127
)

type SchemaMap struct {
	repo         *db.SchemaPartRepository
	schema       *types.Schema
	area         *mt.Area
	cache        *expirable.LRU[int64, *parser.ParsedSchemaPart]
	changedparts map[int64]*parser.ParsedSchemaPart
}

func NewSchemaMap(repo *db.SchemaPartRepository, schema *types.Schema) *SchemaMap {
	return &SchemaMap{
		repo:         repo,
		schema:       schema,
		area:         mt.NewArea(&mt.PosZero, mt.NewPos(schema.SizeX, schema.SizeY, schema.SizeZ).Subtract(mt.NewPos(1, 1, 1))),
		cache:        expirable.NewLRU[int64, *parser.ParsedSchemaPart](1000, nil, time.Second*10),
		changedparts: make(map[int64]*parser.ParsedSchemaPart),
	}
}

func (m *SchemaMap) getPart(mbpos *mt.Pos) (*parser.ParsedSchemaPart, error) {
	key := MapKey(mbpos)

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

func (m *SchemaMap) createBlock(mbpos *mt.Pos) *parser.ParsedSchemaPart {
	pos1 := mbpos.Multiply(16)
	pos2 := pos1.Add(mt.NewPos(15, 15, 15))
	// clip to schematic size
	area := mt.NewArea(pos1, pos2).Union(m.area)
	if area == nil {
		// outside of schematic size
		return nil
	}
	size := area.Size()

	psp := &parser.ParsedSchemaPart{
		PosX:    mbpos.X(),
		PosY:    mbpos.Y(),
		PosZ:    mbpos.Z(),
		NodeIDS: make([]int16, area.Volume()),
		Param1:  make([]byte, area.Volume()),
		Param2:  make([]byte, area.Volume()),
		Meta: &parser.SchemaPartMetadata{
			NodeMapping: map[string]int16{},
			Size: parser.SchemaPartSize{
				X: size.X(),
				Y: size.Y(),
				Z: size.Z(),
			},
			Metadata: &parser.Metadata{
				Meta:   map[string]*parser.MetadataEntry{},
				Timers: map[string]*parser.MetadataTimer{},
			},
		},
	}

	key := MapKey(mbpos)
	m.changedparts[key] = psp

	return psp
}

func (m *SchemaMap) SetNode(pos *mt.Pos, node *mt.Node, md *parser.MetadataEntry) error {
	mbpos := pos.Divide(16).Multiply(16)

	mapblock, err := m.getPart(mbpos)
	if err != nil {
		return fmt.Errorf("getpart error: %v", err)
	}
	if mapblock == nil {
		// create new mapblock
		mapblock = m.createBlock(mbpos)
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

	// node data
	mapblock.Param1[index] = byte(node.Param1)
	mapblock.Param2[index] = byte(node.Param2)
	mapblock.NodeIDS[index] = nodeid

	// metadata
	meta := mapblock.Meta.Metadata
	meta_key := meta.GetKey(rel_pos.X(), rel_pos.Y(), rel_pos.Z())

	if md != nil {
		// overwrite
		entry := meta.Meta[meta_key]
		entry.Fields = md.Fields
		entry.Inventories = md.Inventories
	} else {
		// clear
		meta.Meta[meta_key] = nil
	}

	// mark changed
	key := MapKey(mbpos)
	m.changedparts[key] = mapblock

	return nil
}

func (m *SchemaMap) GetNode(pos *mt.Pos) (*mt.Node, error) {
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

func (m *SchemaMap) Close() error {

	for key, block := range m.changedparts {
		pos := ParseMapKey(key)

		part, err := block.Convert()
		if err != nil {
			return fmt.Errorf("error converting part @ %v: %v", pos, err)
		}

		err = m.repo.CreateOrUpdateSchemaPart(part)
		if err != nil {
			return fmt.Errorf("error updating schemapart @ %v: %v", pos, err)
		}
	}

	return nil
}
