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
		SizeX: 16,
		SizeY: 16,
		SizeZ: 16,
	}

	mods := []types.SchemaMod{
		{ModName: "blah"},
	}

	buf := bytes.NewBuffer([]byte{})
	err := ExportBXSchema(buf, &schema, mods, it)
	assert.NoError(t, err)
	assert.True(t, buf.Len() > 0)
	fmt.Println(buf.String())

	/*
		f, err := os.Create("my.zip")
		assert.NoError(t, err)
		io.Copy(f, bytes.NewReader(buf.Bytes()))
	*/
}
