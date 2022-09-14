package worldedit_test

import (
	"blockexchange/worldedit"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	data, err := os.ReadFile("testdata/plain_chest.we")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	entries, err := worldedit.Parse(data)
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

	fields := part.Meta.Metadata.Meta.Fields["(0,0,0)"]
	assert.Equal(t, "Naj", fields["owner"])

	inventory := part.Meta.Metadata.Meta.Inventories["(0,0,0)"]
	assert.NotNil(t, inventory["main"])
	assert.Equal(t, 32, len(inventory["main"]))
	assert.Equal(t, "default:stone 10", inventory["main"][0])
}

func TestImportComplex(t *testing.T) {
	data, err := os.ReadFile("testdata/zlard.we")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	entries, err := worldedit.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, entries)

	parts, err := worldedit.Import(entries)
	assert.NoError(t, err)
	assert.NotNil(t, parts)

	assert.Equal(t, 8, len(parts))
}
