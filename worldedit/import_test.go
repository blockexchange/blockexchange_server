package worldedit_test

import (
	"blockexchange/worldedit"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	data, err := os.ReadFile("testdata/plain_chest.we")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	entries, err := worldedit.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, entries)

	parts, err := worldedit.Import(entries)
	assert.NoError(t, err)
	assert.NotNil(t, parts)

	//TODO: validate content
}
