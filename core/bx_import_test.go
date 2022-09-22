package core

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBXImport(t *testing.T) {
	f, err := os.Open("testdata/blockexchange.zip")
	assert.NoError(t, err)
	assert.NotNil(t, f)

	stat, err := f.Stat()
	assert.NoError(t, err)

	res, err := ImportBXSchema(f, stat.Size())
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
