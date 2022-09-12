package importwip_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

func TestWEImport(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	data, err := os.ReadFile("plain_chest.we")
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, data[0], uint8('5'))
	assert.Equal(t, data[1], uint8(':'))

	fn, err := L.LoadString(string(data[2:]))
	assert.NoError(t, err)
	assert.NotNil(t, fn)
}
