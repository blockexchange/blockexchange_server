package parser

import (
	"blockexchange/types"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntConversion(t *testing.T) {
	offset := 32768
	high := 0x80
	low := 0x03

	value := int16((int(high) * 256) + int(low) - offset)
	assert.Equal(t, int16(3), value)
}

func TestParseSchemaPart(t *testing.T) {
	f, err := os.Open("testdata/schemapart_1.json")
	assert.NoError(t, err)
	assert.NotNil(t, f)

	part := types.SchemaPart{}
	err = json.NewDecoder(f).Decode(&part)
	assert.NoError(t, err)
	assert.NotNil(t, f)

	parsed, err := ParseSchemaPart(&part)
	assert.NoError(t, err)
	assert.NotNil(t, parsed)
	assert.NotNil(t, parsed.Meta)
	assert.NotNil(t, parsed.Meta.NodeMapping)
	assert.NotNil(t, parsed.Meta.Size)
	assert.NotNil(t, parsed.NodeIDS)
	assert.NotNil(t, parsed.Param1)
	assert.NotNil(t, parsed.Param2)

	id_node_mapping := make(map[int16]string)
	for k, v := range parsed.Meta.NodeMapping {
		id_node_mapping[int16(v)] = k
	}
	for i, v := range parsed.NodeIDS {
		nodename := id_node_mapping[v]
		if nodename == "" {
			t.Errorf("node_id not found: %d @ position %d", v, i)
		}
	}

	for i, v := range parsed.Param1 {
		if v > 15 {
			t.Errorf("param1 too high: %d @ position %d", v, i)
		}
	}
}
