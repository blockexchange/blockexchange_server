package worldedit_test

import (
	"blockexchange/types"
	"blockexchange/worldedit"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestSchemapartIterator() func() (*types.SchemaPart, error) {
	i := 0
	return func() (*types.SchemaPart, error) {
		if i == 0 {
			i++
			f, err := os.Open("testdata/metadata-block/schemapart_0_0_0.json")
			if err != nil {
				return nil, err
			}
			part := types.SchemaPart{}
			err = json.NewDecoder(f).Decode(&part)
			return &part, err
		}

		return nil, nil
	}
}

func TestExportMetadata(t *testing.T) {
	it := createTestSchemapartIterator()
	buf := bytes.NewBuffer([]byte{})
	err := worldedit.Export(buf, it)
	assert.NoError(t, err)
	assert.True(t, buf.Len() > 0)
}

func TestExportImportRoundtrip(t *testing.T) {
	// from blockexchange schema
	it := createTestSchemapartIterator()
	buf := bytes.NewBuffer([]byte{})

	// to worldedit schema
	err := worldedit.Export(buf, it)
	assert.NoError(t, err)
	assert.True(t, buf.Len() > 0)

	// to WE-entries
	entries, modnames, err := worldedit.Parse(buf.Bytes())
	assert.NoError(t, err)
	assert.NotNil(t, entries)
	assert.NotNil(t, modnames)
	assert.Equal(t, 261, len(entries))
	assert.Equal(t, 3, len(modnames))

	// to parsed schemapart
	parts, err := worldedit.Import(entries)
	assert.NoError(t, err)
	assert.NotNil(t, parts)
	assert.Equal(t, 1, len(parts))

	part := parts[0]
	assert.Equal(t, 0, part.PosX)
	assert.Equal(t, 0, part.PosY)
	assert.Equal(t, 0, part.PosZ)
	assert.Equal(t, 16, part.Meta.Size.X)
	assert.Equal(t, 2, part.Meta.Size.Y) // full block size info got lost in we-schema
	assert.Equal(t, 16, part.Meta.Size.Z)

	luac_key := part.Meta.Metadata.GetKey(5, 1, 1)
	luac_entry := part.Meta.Metadata.Meta[luac_key]
	assert.NotNil(t, luac_entry)
	assert.Equal(t, "-- simple content", luac_entry.Fields["code"])

	shelf_key := part.Meta.Metadata.GetKey(3, 1, 1)
	shelf_entry := part.Meta.Metadata.Meta[shelf_key]
	assert.NotNil(t, shelf_entry)
	assert.Equal(t, 1, len(shelf_entry.Inventories))
	assert.Equal(t, 16, len(shelf_entry.Inventories["books"]))
}

func TestExportSchemaSimple(t *testing.T) {

	i := 0
	it := func() (*types.SchemaPart, error) {
		if i >= 4 {
			return nil, nil
		}
		i++

		f, err := os.Open(fmt.Sprintf("testdata/schemapart_%d.json", i))
		if err != nil {
			return nil, err
		}
		part := types.SchemaPart{}
		err = json.NewDecoder(f).Decode(&part)
		return &part, err
	}

	buf := bytes.NewBuffer([]byte{})
	err := worldedit.Export(buf, it)
	assert.NoError(t, err)
	assert.True(t, buf.Len() > 0)
}
