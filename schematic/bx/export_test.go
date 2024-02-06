package bx_test

import (
	"blockexchange/schematic/bx"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func get_schemapart(t *testing.T, filename string) *types.SchemaPart {
	f, err := os.Open(filename)
	assert.NoError(t, err)
	part := &types.SchemaPart{}
	err = json.NewDecoder(f).Decode(part)
	assert.NoError(t, err)
	return part
}

func TestExportBXSchema(t *testing.T) {
	schema := types.Schema{
		Name:  "",
		SizeX: 32,
		SizeY: 16,
		SizeZ: 32,
	}

	mods := []*types.SchemaMod{
		{ModName: "blah"},
	}

	buf := bytes.NewBuffer([]byte{})
	e := bx.NewExporter(buf)
	assert.NoError(t, e.ExportMetadata(&schema, mods))

	for i := 1; i <= 4; i++ {
		err := e.Export(get_schemapart(t, fmt.Sprintf("testdata/schemapart_%d.json", i)))
		assert.NoError(t, err)
	}

	assert.NoError(t, e.Close())
	assert.True(t, buf.Len() > 0)

	// import

	res, err := bx.ParseBXContent(buf.Bytes())
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Mods)
	assert.NotNil(t, res.Parts)
	assert.NotNil(t, res.Schema)

	assert.Equal(t, 1, len(res.Mods))
	assert.Equal(t, "blah", res.Mods[0])

	assert.Equal(t, 32, res.Schema.SizeX)
	assert.Equal(t, 16, res.Schema.SizeY)
	assert.Equal(t, 32, res.Schema.SizeZ)

	assert.Equal(t, 4, len(res.Parts))
}
