package worldedit_test

import (
	"blockexchange/worldedit"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseWorldedit(t *testing.T) {
	data, err := os.ReadFile("testdata/plain_chest.we")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	entries, err := worldedit.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, entries)

	fmt.Println(entries)
}
