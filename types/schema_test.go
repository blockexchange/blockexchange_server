package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaMarshalling(t *testing.T) {
	schema := Schema{
		ID:   1,
		Name: "blah",
	}

	data, err := json.Marshal(schema)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	schema2 := Schema{}
	err = json.Unmarshal(data, &schema2)
	assert.NoError(t, err)
	assert.Equal(t, schema.ID, schema2.ID)
}
