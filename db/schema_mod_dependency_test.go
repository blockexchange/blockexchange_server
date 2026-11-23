package db_test

import (
	"blockexchange/testutils"
	"blockexchange/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaModDependency(t *testing.T) {
	repos := CreateTestDatabase(t)
	user := testutils.CreateUser(repos.UserRepo, t, nil)
	schema := testutils.CreateSchema(repos.SchemaRepo, t, user, nil)

	list, err := repos.SchemaModDependencyRepo.GetBySchemaUID(schema.UID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(list))

	// create, version defaults to "latest"
	game := &types.SchemaModDependency{
		SchemaUID: schema.UID,
		ModType:   types.SchemaModDependencyTypeGame,
		GitURL:    "https://github.com/luanti-org/minetest_game",
	}
	assert.NoError(t, repos.SchemaModDependencyRepo.Create(game))
	assert.NotEmpty(t, game.UID)

	// multiple dependencies per schema
	mod := &types.SchemaModDependency{
		SchemaUID: schema.UID,
		ModType:   types.SchemaModDependencyTypeMod,
		GitURL:    "https://github.com/example/mod",
		Version:   "v1.0.0",
	}
	assert.NoError(t, repos.SchemaModDependencyRepo.Create(mod))

	list, err = repos.SchemaModDependencyRepo.GetBySchemaUID(schema.UID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(list))

	d, err := repos.SchemaModDependencyRepo.GetByUID(game.UID)
	assert.NoError(t, err)
	assert.NotNil(t, d)
	assert.Equal(t, types.SchemaModDependencyTypeGame, d.ModType)
	assert.Equal(t, "latest", d.Version)

	// update
	d.Version = "v1.2.3"
	assert.NoError(t, repos.SchemaModDependencyRepo.Update(d))
	d, err = repos.SchemaModDependencyRepo.GetByUID(game.UID)
	assert.NoError(t, err)
	assert.Equal(t, "v1.2.3", d.Version)

	// invalid mod_type is rejected by the db
	d.ModType = "other"
	assert.Error(t, repos.SchemaModDependencyRepo.Update(d))

	// delete
	assert.NoError(t, repos.SchemaModDependencyRepo.Delete(game.UID))
	d, err = repos.SchemaModDependencyRepo.GetByUID(game.UID)
	assert.NoError(t, err)
	assert.Nil(t, d)

	list, err = repos.SchemaModDependencyRepo.GetBySchemaUID(schema.UID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))
}
