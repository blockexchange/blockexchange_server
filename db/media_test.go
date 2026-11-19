package db_test

import (
	"blockexchange/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMedia(t *testing.T) {
	repos := CreateTestDatabase(t)
	mr := repos.MediaRepo

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
}
