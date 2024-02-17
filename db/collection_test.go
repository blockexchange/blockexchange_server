package db_test

import (
	"blockexchange/db"
	"blockexchange/types"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollection(t *testing.T) {
	kdb := CreateTestDatabase(t)
	repos := db.NewRepositories(kdb)

	u := &types.User{
		Name: fmt.Sprintf("test_%d", rand.Intn(1000)),
	}

	// create

	assert.NoError(t, repos.UserRepo.CreateUser(u))
	c := &types.Collection{
		UserUID: u.UID,
		Name:    "test",
	}
	assert.NoError(t, repos.CollectionRepo.CreateCollection(c))

	// read

	list, err := repos.CollectionRepo.GetCollectionsByUserUID(u.UID)
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "test", list[0].Name)

	// update

	c.Name = "test2"
	assert.NoError(t, repos.CollectionRepo.UpdateCollection(c))

	c2, err := repos.CollectionRepo.GetCollectionByUserUIDAndName(u.UID, "test2")
	assert.NoError(t, err)
	assert.NotNil(t, c2)
	assert.Equal(t, "test2", c2.Name)
	assert.Equal(t, c.UID, c2.UID)

	// delete

	assert.NoError(t, repos.CollectionRepo.DeleteCollection(c.UID))
}
