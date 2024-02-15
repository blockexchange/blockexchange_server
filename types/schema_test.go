package types

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSchemaMarshalling(t *testing.T) {
	schema := Schema{
		UID:  uuid.NewString(),
		Name: "blah",
	}

	data, err := json.Marshal(schema)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	schema2 := Schema{}
	err = json.Unmarshal(data, &schema2)
	assert.NoError(t, err)
	assert.Equal(t, schema.UID, schema2.UID)
}
