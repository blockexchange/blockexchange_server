package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetColorMapping(t *testing.T) {
	m, err := GetColorMapping()
	assert.NoError(t, err)
	assert.NotNil(t, m)
	assert.NotNil(t, m["default:stone"])
}
