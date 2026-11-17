package worldedit_test

import (
	"blockexchange/schematic/worldedit"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	data, err := os.ReadFile("testdata/plain_chest.we")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	entries, _, err := worldedit.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, entries)

	parts, err := worldedit.Import(entries)
	assert.NoError(t, err)
	assert.NotNil(t, parts)

	assert.Equal(t, 1, len(parts))
	part := parts[0]
	assert.Equal(t, 0, part.PosX)
	assert.Equal(t, 0, part.PosY)
	assert.Equal(t, 0, part.PosZ)
	assert.Equal(t, 1, part.Meta.Size.X)
	assert.Equal(t, 1, part.Meta.Size.Y)
	assert.Equal(t, 1, part.Meta.Size.Z)

	assert.Equal(t, 1, len(part.NodeIDS))
	assert.Equal(t, 1, len(part.Param1))
	assert.Equal(t, 1, len(part.Param2))

	assert.Equal(t, int16(1), part.NodeIDS[0])
	assert.Equal(t, byte(140), part.Param1[0])
	assert.Equal(t, byte(3), part.Param2[0])

	assert.NotNil(t, part.Meta.Metadata.Meta)

	fields := part.Meta.Metadata.Meta["(0,0,0)"].Fields
	assert.Equal(t, "Naj", fields["owner"])

	inventory := part.Meta.Metadata.Meta["(0,0,0)"].Inventories
	assert.NotNil(t, inventory["main"])
	assert.Equal(t, 32, len(inventory["main"]))
	assert.Equal(t, "default:stone 10", inventory["main"][0])
}

func TestImportComplex(t *testing.T) {
	data, err := os.ReadFile("testdata/zlard.we")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	entries, _, err := worldedit.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, entries)

	parts, err := worldedit.Import(entries)
	assert.NoError(t, err)
	assert.NotNil(t, parts)

	assert.Equal(t, 8, len(parts))
}

func TestImportNodeIDsPerPart(t *testing.T) {
	// 40000 distinct node names would overflow a global int16 node-id counter
	entries := make([]*worldedit.WEEntry, 0)
	for i := 0; i < 40000; i++ {
		entries = append(entries, &worldedit.WEEntry{
			Name: fmt.Sprintf("mod:node_%d", i),
			X:    i % 200,
			Y:    (i / 200) % 200,
			Z:    0,
		})
	}

	parts, err := worldedit.Import(entries)
	assert.NoError(t, err)
	assert.NotEmpty(t, parts)

	for _, part := range parts {
		mapping := part.Meta.NodeMapping
		assert.Equal(t, int16(0), mapping["air"])

		// ids are dense in [0, len(mapping))
		seen := map[int16]bool{}
		for _, id := range mapping {
			assert.True(t, id >= 0 && int(id) < len(mapping), "id %d out of range", id)
			assert.False(t, seen[id], "duplicate id %d", id)
			seen[id] = true
		}
		for _, id := range part.NodeIDS {
			assert.True(t, id >= 0 && int(id) < len(mapping))
		}
	}
}
