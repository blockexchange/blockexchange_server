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

func TestModMedia(t *testing.T) {
	repos := CreateTestDatabase(t)
	mr := repos.MediaRepo

	mod := &types.Mod{
		Name:         "mymod",
		Source:       "https://...",
		MediaLicense: "CC-BY-3.0",
	}

	// prepare
	mr.RemoveMod("mymod")
	assert.NoError(t, mr.CreateMod(mod))

	// create nodedef
	nd := &types.Nodedefinition{
		Name:       "mymod:xy",
		ModName:    "mymod",
		Definition: "{}",
	}
	assert.NoError(t, mr.CreateNodedefinition(nd))

	// get nodedef
	nd2, err := mr.GetNodedefinitionByName("mymod:xy")
	assert.NoError(t, err)
	assert.NotNil(t, nd2)
	assert.EqualValues(t, nd, nd2)

	// get all nodedefs
	nodedefs, err := mr.GetNodedefinitions()
	assert.NoError(t, err)
	assert.NotNil(t, nodedefs)
	assert.True(t, len(nodedefs) >= 1)

	// update
	nd.Definition = `{"x":1}`
	assert.NoError(t, mr.UpdateNodedefinition(nd))

	// remove
	assert.NoError(t, mr.RemoveNodedefinition("mymod:xy"))

	// create media
	mf := &types.Mediafile{
		Name:    "default_stone.png",
		ModName: "mymod",
		Data:    []byte{0x00},
	}
	assert.NoError(t, mr.CreateMediafile(mf))

	// get media
	mf2, err := mr.GetMediafileByName("default_stone.png")
	assert.NoError(t, err)
	assert.NotNil(t, mf2)
	assert.EqualValues(t, mf, mf2)

	// remove media
	assert.NoError(t, mr.RemoveMediafile("default_stone.png"))

	// remove all
	assert.NoError(t, mr.RemoveMod("mymod"))

}
