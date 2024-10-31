package db_test

import (
	"blockexchange/types"
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetSchemaTagsByIDs(t *testing.T) {
	repos := CreateTestDatabase(t)

	u := &types.User{
		Name: fmt.Sprintf("test_%d", rand.Intn(10000)),
	}
	assert.NoError(t, repos.UserRepo.CreateUser(u))
	s := &types.Schema{
		UID:     uuid.NewString(),
		UserUID: u.UID,
		Name:    "test",
	}
	assert.NoError(t, repos.SchemaRepo.CreateSchema(s))

	t1 := &types.Tag{Name: fmt.Sprintf("tag_%d", rand.Intn(10000))}
	assert.NoError(t, repos.TagRepo.Create(t1))

	t2 := &types.Tag{Name: fmt.Sprintf("tag_%d", rand.Intn(10000))}
	assert.NoError(t, repos.TagRepo.Create(t2))

	st1 := &types.SchemaTag{TagUID: t1.UID, SchemaUID: s.UID}
	assert.NoError(t, repos.SchemaTagRepo.Create(st1))

	st2 := &types.SchemaTag{TagUID: t1.UID, SchemaUID: s.UID}
	assert.NoError(t, repos.SchemaTagRepo.Create(st2))
}
