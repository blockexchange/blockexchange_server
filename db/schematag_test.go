package db_test

import (
	"blockexchange/db"
	"blockexchange/types"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSchemaTagsByIDs(t *testing.T) {
	kdb := CreateTestDatabase(t)
	repos := db.NewRepositories(kdb)

	u := &types.User{
		Name: fmt.Sprintf("test_%d", rand.Intn(10000)),
	}
	assert.NoError(t, repos.UserRepo.CreateUser(u))
	s := &types.Schema{
		UserID: *u.ID,
		Name:   "test",
	}
	assert.NoError(t, repos.SchemaRepo.CreateSchema(s))

	t1 := &types.Tag{Name: fmt.Sprintf("tag_%d", rand.Intn(10000))}
	assert.NoError(t, repos.TagRepo.Create(t1))

	t2 := &types.Tag{Name: fmt.Sprintf("tag_%d", rand.Intn(10000))}
	assert.NoError(t, repos.TagRepo.Create(t2))

	st1 := &types.SchemaTag{TagUID: t1.UID, SchemaID: *s.ID}
	assert.NoError(t, repos.SchemaTagRepo.Create(st1))

	st2 := &types.SchemaTag{TagUID: t1.UID, SchemaID: *s.ID}
	assert.NoError(t, repos.SchemaTagRepo.Create(st2))

	sts, err := repos.SchemaTagRepo.GetBySchemaIDs([]int64{*s.ID})
	assert.NoError(t, err)
	assert.NotNil(t, sts)
	assert.Equal(t, 2, len(sts))
}
