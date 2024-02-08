package db_test

import (
	"blockexchange/db"
	"blockexchange/types"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSchemaModsByIDs(t *testing.T) {
	kdb := CreateTestDatabase(t)
	repos := db.NewRepositories(kdb)

	u := &types.User{
		Name: fmt.Sprintf("test_%d", rand.Intn(1000)),
	}
	assert.NoError(t, repos.UserRepo.CreateUser(u))
	s := &types.Schema{
		UserID: *u.ID,
		Name:   "test",
	}
	assert.NoError(t, repos.SchemaRepo.CreateSchema(s))

	assert.NoError(t, repos.SchemaModRepo.CreateSchemaMod(&types.SchemaMod{
		SchemaID: *s.ID,
		ModName:  "mod1",
	}))
	assert.NoError(t, repos.SchemaModRepo.CreateSchemaMod(&types.SchemaMod{
		SchemaID: *s.ID,
		ModName:  "mod2",
	}))

	sm, err := repos.SchemaModRepo.GetSchemaModsBySchemaIDs([]int64{*s.ID})
	assert.NoError(t, err)
	assert.NotNil(t, sm)
	assert.Equal(t, 2, len(sm))
}
