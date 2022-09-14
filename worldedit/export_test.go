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

func TestExportSchemaMetadata(t *testing.T) {

	i := 0
	it := func() (*types.SchemaPart, error) {
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

	buf := bytes.NewBuffer([]byte{})
	err := worldedit.Export(buf, it)
	assert.NoError(t, err)
	assert.True(t, buf.Len() > 0)
	fmt.Println(buf.String())
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
	fmt.Println(buf.String())
}
