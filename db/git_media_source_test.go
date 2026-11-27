package db_test

import (
	"blockexchange/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitMediaSource(t *testing.T) {
	repos := CreateTestDatabase(t)
	r := repos.GitMediaSourceRepo

	s := &types.GitMediaSource{GitURL: "https://github.com/mt-mods/gitsource-test"}
	assert.NoError(t, r.Create(s))
	assert.NotEmpty(t, s.ID)

	// the git url is unique
	assert.Error(t, r.Create(&types.GitMediaSource{GitURL: s.GitURL}))

	s.License = "MIT"
	assert.NoError(t, r.Update(s))
	assert.NoError(t, r.UpdateLicenses(s.ID, "mt-mods/gitsource-test", "CC0-1.0", "MIT"))
	assert.NoError(t, r.UpdateStatus(s.ID, 1234, true, "ok"))

	s2, err := r.GetByID(s.ID)
	assert.NoError(t, err)
	assert.NotNil(t, s2)
	assert.Equal(t, s.GitURL, s2.GitURL)
	assert.Equal(t, "mt-mods/gitsource-test", s2.CDBPackage)
	assert.Equal(t, "CC0-1.0", s2.MediaLicense)
	assert.Equal(t, "MIT", s2.License)
	assert.Equal(t, int64(1234), s2.LastUpdatedTimestamp)
	assert.True(t, s2.LastUpdatedSuccess)
	assert.Equal(t, "ok", s2.LastUpdatedMessage)

	list, err := r.GetAll()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)

	assert.NoError(t, r.Delete(s.ID))
	s2, err = r.GetByID(s.ID)
	assert.NoError(t, err)
	assert.Nil(t, s2)
}
