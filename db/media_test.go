package db_test

import (
	"blockexchange/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModEntity(t *testing.T) {
	repos := CreateTestDatabase(t)
	mr := repos.MediaRepo

	mod := &types.Mod{
		Name:         "mymod",
		Source:       "https://...",
		CodeLicense:  "MIT",
		MediaLicense: "CC-BY-3.0",
	}

	// cleanup
	mr.RemoveMod("mymod")

	// create
	assert.NoError(t, mr.CreateMod(mod))

	// read
	mod2, err := mr.GetModByName("mymod")
	assert.NoError(t, err)
	assert.NotNil(t, mod2)
	assert.EqualValues(t, mod, mod2)

	// update
	mod.Source = "http://..."
	assert.NoError(t, mr.UpdateMod(mod))

	// read
	mod2, err = mr.GetModByName("mymod")
	assert.NoError(t, err)
	assert.NotNil(t, mod2)
	assert.EqualValues(t, mod, mod2)

	// delete
	assert.NoError(t, mr.RemoveMod("mymod"))

	// read
	mod2, err = mr.GetModByName("mymod")
	assert.NoError(t, err)
	assert.Nil(t, mod2)
}
