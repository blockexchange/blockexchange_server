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
		Name: "default_stone.png",
		Data: []byte{0x00},
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
	assert.Nil(t, infos[0].GitMediaSourceID)
	assert.Nil(t, infos[0].GitURL)
	assert.Equal(t, int64(1), infos[0].Size)
	c, err = mr.CountMediafiles(db.MediaQuery{Query: "default_stone.png"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), c)

	// remove media
	assert.NoError(t, mr.RemoveMediafile("default_stone.png"))
}

func TestMediafileGitSource(t *testing.T) {
	repos := CreateTestDatabase(t)
	mr := repos.MediaRepo

	src := &types.GitMediaSource{GitURL: "https://github.com/test/mediasource-test", MediaLicense: "CC-BY-SA-4.0"}
	assert.NoError(t, repos.GitMediaSourceRepo.Create(src))
	defer repos.GitMediaSourceRepo.Delete(src.ID)

	assert.NoError(t, mr.UpsertMediafiles([]*types.Mediafile{
		{Name: "mediasource_a.png", GitMediaSourceID: &src.ID, Data: []byte{1}},
		{Name: "mediasource_b.png", GitMediaSourceID: &src.ID, Data: []byte{2, 3}},
	}))

	infos, err := mr.SearchMediafileInfos(db.MediaQuery{Query: "mediasource_", Limit: -1})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(infos))
	assert.Equal(t, src.ID, *infos[0].GitMediaSourceID)
	assert.Equal(t, src.GitURL, *infos[0].GitURL)
	assert.Equal(t, "CC-BY-SA-4.0", *infos[0].MediaLicense)
	assert.Equal(t, int64(2), infos[1].Size)

	// license filter
	unlicensed := &types.GitMediaSource{GitURL: "https://github.com/test/mediasource-nolicense"}
	assert.NoError(t, repos.GitMediaSourceRepo.Create(unlicensed))
	defer repos.GitMediaSourceRepo.Delete(unlicensed.ID)
	assert.NoError(t, mr.UpsertMediafiles([]*types.Mediafile{
		{Name: "mediasource_c.png", GitMediaSourceID: &unlicensed.ID, Data: []byte{4}},
		{Name: "mediasource_d.png", Data: []byte{5}},
	}))
	defer mr.RemoveMediafile("mediasource_d.png")

	with, without := true, false
	infos, err = mr.SearchMediafileInfos(db.MediaQuery{Query: "mediasource_", License: &with, Limit: -1})
	assert.NoError(t, err)
	assert.Equal(t, []string{"mediasource_a.png", "mediasource_b.png"}, []string{infos[0].Name, infos[1].Name})
	assert.Equal(t, 2, len(infos))
	infos, err = mr.SearchMediafileInfos(db.MediaQuery{Query: "mediasource_", License: &without, Limit: -1})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(infos))
	assert.Equal(t, "mediasource_c.png", infos[0].Name)
	assert.Equal(t, "mediasource_d.png", infos[1].Name)
	c, err := mr.CountMediafiles(db.MediaQuery{Query: "mediasource_", License: &with})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), c)
	c, err = mr.CountMediafiles(db.MediaQuery{Query: "mediasource_", License: &without})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), c)
	c, err = mr.CountMediafiles(db.MediaQuery{Query: "mediasource_"})
	assert.NoError(t, err)
	assert.Equal(t, int64(4), c)

	names, err := mr.GetMediafileNamesBySource(src.ID)
	assert.NoError(t, err)
	assert.Equal(t, []string{"mediasource_a.png", "mediasource_b.png"}, names)

	// only the files of the source are removed
	assert.NoError(t, mr.RemoveMediafilesOfSource("00000000-0000-0000-0000-000000000000", []string{"mediasource_a.png"}))
	assert.NoError(t, mr.RemoveMediafilesOfSource(src.ID, []string{"mediasource_a.png"}))
	names, err = mr.GetMediafileNamesBySource(src.ID)
	assert.NoError(t, err)
	assert.Equal(t, []string{"mediasource_b.png"}, names)

	// removing the source removes its files
	assert.NoError(t, repos.GitMediaSourceRepo.Delete(src.ID))
	mf, err := mr.GetMediafileByName("mediasource_b.png")
	assert.NoError(t, err)
	assert.Nil(t, mf)
}
