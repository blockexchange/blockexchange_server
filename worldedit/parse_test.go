package worldedit_test

import (
	"blockexchange/worldedit"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseWorldedit(t *testing.T) {
	data, err := os.ReadFile("testdata/plain_chest.we")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	entries, err := worldedit.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, entries)

	fmt.Println(entries)

	assert.Equal(t, 1, len(entries))
	e := entries[0]
	assert.Equal(t, "default:chest_locked", e.Name)
	assert.Equal(t, 0, e.X)
	assert.Equal(t, 0, e.Y)
	assert.Equal(t, 0, e.Z)
	assert.Equal(t, 140, e.Param1)
	assert.Equal(t, 3, e.Param2)

	assert.NotNil(t, e.Meta)
	assert.NotNil(t, e.Meta.Fields)
	assert.Equal(t, "Naj", e.Meta.Fields["owner"])

	assert.NotNil(t, e.Meta.Inventory)
	assert.NotNil(t, e.Meta.Inventory["main"])
	assert.Equal(t, 32, len(e.Meta.Inventory["main"]))
	assert.Equal(t, "default:stone 10", e.Meta.Inventory["main"][0])
}
