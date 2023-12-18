package core

import (
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExportBXSchema(t *testing.T) {

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

	schema := types.Schema{
		Name:  "",
		SizeX: 32,
		SizeY: 16,
		SizeZ: 32,
	}

	mods := []types.SchemaMod{
		{ModName: "blah"},
	}

	buf := bytes.NewBuffer([]byte{})
	err := ExportBXSchema(buf, &schema, mods, it)
	assert.NoError(t, err)
	assert.True(t, buf.Len() > 0)

	// import

	res, err := ParseBXContent(buf.Bytes())
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

	/*
		f, err := os.Create("my.zip")
		assert.NoError(t, err)
		io.Copy(f, bytes.NewReader(buf.Bytes()))
	*/
}
