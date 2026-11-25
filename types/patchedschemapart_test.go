package types

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPatchedSchemaPartSerialization(t *testing.T) {
	part := PatchedSchemapart{
		SchemaPart: &SchemaPart{
			SchemaUID: uuid.NewString(),
			OffsetX:   3,
			OffsetY:   4,
			OffsetZ:   5,
			Mtime:     6,
			Data:      []byte{0x07, 0x08},
			MetaData:  []byte{0x09, 0x0A, 0x0B},
		},
		PatchFilter: []PatchFilterType{PatchNoOp, PatchReplaceAll, PatchReplaceParam2},
	}

	data, err := json.Marshal(part)
	assert.NoError(t, err)

	m := map[string]any{}
	assert.NoError(t, json.Unmarshal(data, &m))
	assert.Equal(t, base64.RawStdEncoding.EncodeToString([]byte{0x00, 0x0F, 0x04}), m["patchfilter"])
	assert.Equal(t, base64.RawStdEncoding.EncodeToString(part.Data), m["data"])

	part2 := PatchedSchemapart{}
	assert.NoError(t, json.Unmarshal(data, &part2))
	assert.Equal(t, part.SchemaPart, part2.SchemaPart)
	assert.Equal(t, part.PatchFilter, part2.PatchFilter)
}

func TestPatchedSchemaPartInvalidPatchFilter(t *testing.T) {
	part := PatchedSchemapart{}
	err := json.Unmarshal([]byte(`{"data":"","metadata":"","patchfilter":"!!"}`), &part)
	assert.Error(t, err)
}
