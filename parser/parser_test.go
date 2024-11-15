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

func TestParseSchemaPartPos(t *testing.T) {
	f, err := os.Open("testdata/schemapart_2.json")
	assert.NoError(t, err)
	assert.NotNil(t, f)

	part := &types.SchemaPart{}
	err = json.NewDecoder(f).Decode(part)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	assert.Equal(t, 16, part.OffsetZ)
	assert.Equal(t, 0, part.OffsetY)
	assert.Equal(t, 0, part.OffsetX)

	psp, err := ParseSchemaPart(part)
	assert.NoError(t, err)
	assert.NotNil(t, psp)

	part2, err := psp.Convert()
	assert.NoError(t, err)
	assert.NotNil(t, part2)

	assert.Equal(t, 16, part2.OffsetZ)
	assert.Equal(t, 0, part2.OffsetY)
	assert.Equal(t, 0, part2.OffsetX)
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

	sp2, err := parsed.Convert()
	assert.NoError(t, err)
	assert.NotNil(t, sp2)

	assert.Equal(t, 0, sp2.OffsetX)
	assert.Equal(t, 0, sp2.OffsetY)
	assert.Equal(t, 0, sp2.OffsetZ)

	assert.True(t, len(sp2.Data) > 2)
	assert.True(t, len(sp2.MetaData) > 2)
}
