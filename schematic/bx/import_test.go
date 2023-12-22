package bx_test

import (
	"blockexchange/schematic/bx"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBXImport(t *testing.T) {
	data, err := os.ReadFile("testdata/blockexchange.zip")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	res, err := bx.ParseBXContent(data)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.NotNil(t, res.Mods)
	assert.NotNil(t, res.Schema)
	assert.NotNil(t, res.Parts)

	assert.Equal(t, 86, res.Schema.SizeX)
	assert.Equal(t, 37, res.Schema.SizeY)
	assert.Equal(t, 2, res.Schema.SizeZ)

	assert.Equal(t, 1, len(res.Mods))
	assert.Equal(t, 18, len(res.Parts))
}
