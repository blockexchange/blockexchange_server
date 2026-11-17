package worldedit_test

import (
	"blockexchange/schematic/worldedit"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSyntaxVariants(t *testing.T) {
	data := `5:return {
		{x=1, y=2, z=3, name="default:stone", param1=0x0F, param2=2.0e0},
		{["x"] = 0; ["y"] = 0; ["z"] = 0; ["name"] = 'mod_a:Node_1';},
		{x=0,y=0,z=0,name="air"},
		{x=4,y=5,z=6,name="default:chest", meta={fields={
			esc="a\"b\\c\nd\
e\65\x41\0z\'",
			num=42,
			float=1.5,
		}, inventory={main={"default:dirt 1", [3]="c", [2]="b"}}}},
		{x=0,y=0,z=0,name="default:dirt",meta={fields={},inventory={}}},
	}`

	entries, modnames, err := worldedit.Parse([]byte(data))
	require.NoError(t, err)
	assert.Equal(t, []string{"default", "mod_a"}, modnames)
	// air is skipped
	require.Equal(t, 4, len(entries))

	assert.Equal(t, worldedit.WEEntry{Name: "default:stone", X: 1, Y: 2, Z: 3, Param1: 15, Param2: 2}, *entries[0])
	assert.Equal(t, "mod_a:Node_1", entries[1].Name)

	meta := entries[2].Meta
	require.NotNil(t, meta)
	assert.Equal(t, "a\"b\\c\nd\neAA\x00z'", meta.Fields["esc"])
	assert.Equal(t, "42", meta.Fields["num"])
	assert.Equal(t, "1.5", meta.Fields["float"])
	assert.Equal(t, []string{"default:dirt 1", "b", "c"}, meta.Inventory["main"])

	// empty metadata is dropped
	assert.Nil(t, entries[3].Meta)
}

func TestParseRejectsInvalidInput(t *testing.T) {
	node := func(s string) string { return "5:return {" + s + "}" }

	cases := map[string]string{
		// header
		"empty":          "",
		"one byte":       "5",
		"wrong version":  "4:return {}",
		"missing return": "5:{{x=0,y=0,z=0,name=\"a:b\"}}",
		"empty list":     "5:return {}",
		"trailing data":  node(`{x=0,y=0,z=0,name="a:b"}`) + " x",

		// code execution attempts
		"infinite loop":    "5:while true do end",
		"string doubling":  `5:local s="x" for i=1,40 do s=s..s end return {}`,
		"function call":    node(`{x=0,y=0,z=0,name=f()}`),
		"concat":           node(`{x=0,y=0,z=0,name="a".."b"}`),
		"variable":         node(`{x=y,y=0,z=0,name="a:b"}`),
		"arithmetic":       node(`{x=1+1,y=0,z=0,name="a:b"}`),
		"nan":              node(`{x=0/0,y=0,z=0,name="a:b"}`),
		"math.huge":        node(`{x=math.huge,y=0,z=0,name="a:b"}`),
		"comment":          node(`{x=0,y=0,z=0,name="a:b"} -- hi`),
		"long string":      node(`{x=0,y=0,z=0,name=[[a:b]]}`),
		"function literal": node(`{x=0,y=0,z=0,name=function() end}`),

		// structure
		"entry not table":    node(`1`),
		"keyed root":         "5:return {a={x=0,y=0,z=0,name=\"a:b\"}}",
		"meta not table":     node(`{x=0,y=0,z=0,name="a:b",meta=1}`),
		"fields numeric key": node(`{x=0,y=0,z=0,name="a:b",meta={fields={[1]="x"}}}`),
		"fields bool value":  node(`{x=0,y=0,z=0,name="a:b",meta={fields={a=true}}}`),
		"inventory not list": node(`{x=0,y=0,z=0,name="a:b",meta={inventory={main={a="b"}}}}`),
		"inventory hole":     node(`{x=0,y=0,z=0,name="a:b",meta={inventory={main={[1]="a",[3]="b"}}}}`),
		"inventory number":   node(`{x=0,y=0,z=0,name="a:b",meta={inventory={main={1}}}}`),
		"too deep":           node(`{x=0,y=0,z=0,name="a:b",other=` + strings.Repeat("{", 100) + strings.Repeat("}", 100) + `}`),
		"unterminated":       node(`{x=0,y=0,z=0,name="a:b"`),
		"unterminated str":   node(`{x=0,y=0,z=0,name="a:b}`),
		"newline in string":  node("{x=0,y=0,z=0,name=\"a:\nb\"}"),
		"bad escape":         node(`{x=0,y=0,z=0,name="a:\qb"}`),
		"escape too large":   node(`{x=0,y=0,z=0,name="a:\999"}`),
		"missing separator":  node(`{x=0 y=0,z=0,name="a:b"}`),

		// values
		"missing x":        node(`{y=0,z=0,name="a:b"}`),
		"missing name":     node(`{x=0,y=0,z=0}`),
		"string coord":     node(`{x="0",y=0,z=0,name="a:b"}`),
		"negative coord":   node(`{x=-5,y=0,z=0,name="a:b"}`),
		"huge coord":       node(`{x=1e300,y=0,z=0,name="a:b"}`),
		"coord too large":  node(`{x=500,y=0,z=0,name="a:b"}`),
		"fractional coord": node(`{x=1.5,y=0,z=0,name="a:b"}`),
		"param2 overflow":  node(`{x=0,y=0,z=0,name="a:b",param2=999}`),
		"param1 negative":  node(`{x=0,y=0,z=0,name="a:b",param1=-1}`),
		"empty name":       node(`{x=0,y=0,z=0,name=""}`),
		"name w/o mod":     node(`{x=0,y=0,z=0,name="stone"}`),
		"name w/ spaces":   node(`{x=0,y=0,z=0,name="a:b c"}`),
		"name too long":    node(`{x=0,y=0,z=0,name="a:` + strings.Repeat("b", 300) + `"}`),
		"uppercase mod":    node(`{x=0,y=0,z=0,name="Default:stone"}`),
		"too many stacks":  node(`{x=0,y=0,z=0,name="a:b",meta={inventory={main={` + strings.Repeat(`"",`, worldedit.MaxInventorySize+1) + `}}}}`),
	}

	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			var err error
			assert.NotPanics(t, func() {
				_, _, err = worldedit.Parse([]byte(data))
			})
			assert.ErrorIs(t, err, worldedit.ErrInvalidSchematic)
		})
	}
}

func TestParseErrorPosition(t *testing.T) {
	_, _, err := worldedit.Parse([]byte("5:return {\n{x=0,y=0,z=0,name=\"a:b\"},\n{x=0,y=0,z=0,name=oops}}"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "line 3")
}

func FuzzParse(f *testing.F) {
	for _, file := range []string{"testdata/plain_chest.we", "testdata/metadata-block/metadata-testblock.we"} {
		data, err := os.ReadFile(file)
		require.NoError(f, err)
		f.Add(data)
	}
	f.Add([]byte(`5:return {{x=0,y=0,z=0,name="a:b",meta={fields={a="\1\x41"},inventory={m={[2]="",""}}}}}`))

	f.Fuzz(func(t *testing.T, data []byte) {
		entries, _, err := worldedit.Parse(data)
		if err != nil {
			return
		}
		for _, e := range entries {
			if e.X < 0 || e.X > worldedit.MaxCoord || e.Y < 0 || e.Y > worldedit.MaxCoord || e.Z < 0 || e.Z > worldedit.MaxCoord {
				t.Fatalf("coordinates out of range: %v", e)
			}
		}
	})
}
