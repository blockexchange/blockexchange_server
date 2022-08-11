package types

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanityBase64(t *testing.T) {
	data := "eJztwQENAAAMAqAHekijm8MNyEdVVVVVVVVVHX8AAAAAAAAAwLwCfjrAlw"
	_, err := base64.RawStdEncoding.DecodeString(data)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSchemaPartSerialization(t *testing.T) {
	part := SchemaPart{
		SchemaID: 2,
		OffsetX:  3,
		OffsetY:  4,
		OffsetZ:  5,
		Mtime:    6,
		Data:     []byte{0x07, 0x08},
		MetaData: []byte{0x09, 0x0A, 0x0B},
	}

	data, err := json.Marshal(part)
	if err != nil {
		t.Fatal(err)
	}

	part2 := SchemaPart{}
	err = json.Unmarshal(data, &part2)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, part.SchemaID, part2.SchemaID)
	assert.Equal(t, part.OffsetX, part2.OffsetX)
	assert.Equal(t, part.OffsetY, part2.OffsetY)
	assert.Equal(t, part.OffsetZ, part2.OffsetZ)
	assert.Equal(t, part.OffsetZ, part2.OffsetZ)
	assert.Equal(t, part.Mtime, part2.Mtime)
	assert.Equal(t, part.MetaData, part2.MetaData)
}

func TestSchemaPartInvalidData(t *testing.T) {
	part := SchemaPart{}
	err := part.UnmarshalJSON([]byte{})
	assert.Error(t, err)
}
