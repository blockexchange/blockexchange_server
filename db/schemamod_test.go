package db_test

import (
	"blockexchange/types"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSchemaModsByIDs(t *testing.T) {
	repos := CreateTestDatabase(t)

	u := &types.User{
		Name: fmt.Sprintf("test_%d", rand.Intn(1000)),
	}
	assert.NoError(t, repos.UserRepo.CreateUser(u))
	s := &types.Schema{
		UserUID: u.UID,
		Name:    "test",
	}
	assert.NoError(t, repos.SchemaRepo.CreateSchema(s))

	assert.NoError(t, repos.SchemaModRepo.CreateSchemaMod(&types.SchemaMod{
		SchemaUID: s.UID,
		ModName:   "mod1",
	}))
	assert.NoError(t, repos.SchemaModRepo.CreateSchemaMod(&types.SchemaMod{
		SchemaUID: s.UID,
		ModName:   "mod2",
	}))

	sm, err := repos.SchemaModRepo.GetSchemaModsBySchemaUID(s.UID)
	assert.NoError(t, err)
	assert.NotNil(t, sm)
	assert.Equal(t, 2, len(sm))
}
