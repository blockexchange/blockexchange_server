package worldedit_test

import (
	"blockexchange/worldedit"
	"os"
	"testing"

	"github.com/Shopify/go-lua"
	"github.com/stretchr/testify/assert"
)

func TestLuaParser(t *testing.T) {
	L := lua.NewState()
	err := lua.DoString(L, "return { {x=1}, {x=2} }")
	assert.NoError(t, err)

	// root table
	assert.True(t, L.IsTable(L.Top()))
	L.Length(L.Top())
	assert.True(t, L.IsNumber(L.Top()))
	length, ok := L.ToNumber(L.Top())
	assert.True(t, ok)
	assert.Equal(t, float64(2), length)
	L.Pop(1)

	// get first entry
	L.PushNumber(1)
	L.Table(L.Top() - 1)
	assert.True(t, L.IsTable(L.Top()))

	// { x=1 }
	assert.True(t, L.IsTable(L.Top()))
	L.PushString("x")
	assert.True(t, L.IsString(L.Top()))
	L.Table(L.Top() - 1)
	assert.True(t, L.IsNumber(L.Top()))
	length, ok = L.ToNumber(L.Top())
	assert.True(t, ok)
	assert.Equal(t, float64(1), length)
	L.Pop(2)

	// get second entry
	L.PushNumber(2)
	L.Table(L.Top() - 1)
	assert.True(t, L.IsTable(L.Top()))

	// { x=2 }
	assert.True(t, L.IsTable(L.Top()))
	L.PushString("x")
	assert.True(t, L.IsString(L.Top()))
	L.Table(L.Top() - 1)
	assert.True(t, L.IsNumber(L.Top()))
	length, ok = L.ToNumber(L.Top())
	assert.True(t, ok)
	assert.Equal(t, float64(2), length)
}

func TestParseWorldedit(t *testing.T) {
	data, err := os.ReadFile("testdata/plain_chest.we")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	entries, modnames, err := worldedit.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, modnames)
	assert.NotNil(t, entries)

	assert.Equal(t, 1, len(modnames))
	assert.Equal(t, "default", modnames[0])

	assert.Equal(t, 1, len(entries))
	e := entries[0]
	assert.Equal(t, "default:chest_locked", e.Name)
	assert.Equal(t, 0, e.X)
	assert.Equal(t, 0, e.Y)
	assert.Equal(t, 0, e.Z)
	assert.Equal(t, byte(140), e.Param1)
	assert.Equal(t, byte(3), e.Param2)

	assert.NotNil(t, e.Meta)
	assert.NotNil(t, e.Meta.Fields)
	assert.Equal(t, "Naj", e.Meta.Fields["owner"])

	assert.NotNil(t, e.Meta.Inventory)
	assert.NotNil(t, e.Meta.Inventory["main"])
	assert.Equal(t, 32, len(e.Meta.Inventory["main"]))
	assert.Equal(t, "default:stone 10", e.Meta.Inventory["main"][0])
}

func TestParseWorldedit2(t *testing.T) {
	data, err := os.ReadFile("testdata/zlard.we")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	entries, _, err := worldedit.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, entries)

	assert.Equal(t, 21382, len(entries))
}

func TestParseWorldedit3(t *testing.T) {
	data, err := os.ReadFile("testdata/deathstar.we")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	entries, _, err := worldedit.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, entries)

	assert.Equal(t, 53646, len(entries))
}
