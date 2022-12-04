package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type JsonTest struct {
	MyInt JsonInt `json:"my_int"`
}

func TestJsonNumber(t *testing.T) {
	str := `{"my_int":22.0}`
	o := &JsonTest{}
	assert.NoError(t, json.Unmarshal([]byte(str), o))
	assert.Equal(t, JsonInt(22), o.MyInt)
}
