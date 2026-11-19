package db_test

import (
	"blockexchange/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLuantiMediaConnection(t *testing.T) {
	repos := CreateTestDatabase(t)

	// create
	c := &types.LuantiMediaConnection{
		Hostname:             "localhost",
		Port:                 30000,
		Username:             "user",
		Password:             "pass",
		LastUpdatedTimestamp: 1234,
		LastUpdatedSuccess:   true,
		LastUpdatedMessage:   "ok",
	}
	assert.NoError(t, repos.LuantiMediaConnectionRepo.Create(c))
	assert.NotEmpty(t, c.ID)

	// read
	c2, err := repos.LuantiMediaConnectionRepo.GetByID(c.ID)
	assert.NoError(t, err)
	assert.NotNil(t, c2)
	assert.Equal(t, "localhost", c2.Hostname)
	assert.Equal(t, 30000, c2.Port)
	assert.Equal(t, int64(1234), c2.LastUpdatedTimestamp)
	assert.True(t, c2.LastUpdatedSuccess)
	assert.Equal(t, "ok", c2.LastUpdatedMessage)

	list, err := repos.LuantiMediaConnectionRepo.GetAll()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)

	// update
	c.Port = 30001
	assert.NoError(t, repos.LuantiMediaConnectionRepo.Update(c))
	c2, err = repos.LuantiMediaConnectionRepo.GetByID(c.ID)
	assert.NoError(t, err)
	assert.Equal(t, 30001, c2.Port)

	// delete
	assert.NoError(t, repos.LuantiMediaConnectionRepo.Delete(c.ID))
	c2, err = repos.LuantiMediaConnectionRepo.GetByID(c.ID)
	assert.NoError(t, err)
	assert.Nil(t, c2)
}
