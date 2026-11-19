package db_test

import (
	"blockexchange/db"
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

	// search and count nodedefs
	nodedefs, err = mr.SearchNodedefinitions(db.MediaQuery{Query: "MYMOD:x", Limit: 10})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(nodedefs))
	c, err := mr.CountNodedefinitions("mymod:x")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), c)
	nodedefs, err = mr.SearchNodedefinitions(db.MediaQuery{Query: "mymod:x", Limit: 10, Offset: 1})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(nodedefs))
	c, err = mr.CountNodedefinitions("mymod%x")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), c)

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

	// list media infos
	infos, err := mr.SearchMediafileInfos(db.MediaQuery{Query: "default_stone.png", Limit: -1})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(infos))
	assert.Equal(t, "mymod", infos[0].ModName)
	assert.Equal(t, int64(1), infos[0].Size)
	c, err = mr.CountMediafiles("default_stone.png")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), c)

	// remove media
	assert.NoError(t, mr.RemoveMediafile("default_stone.png"))
}
