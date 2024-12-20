package core_test

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"blockexchange/types"
	"testing"

	mt "github.com/minetest-go/types"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSchemaMap(t *testing.T) {
	repos := testutils.CreateTestDatabase(t)

	user := &types.User{
		UID:  uuid.NewString(),
		Name: core.CreateToken(10),
		Type: types.UserTypeLocal,
		Role: types.UserRoleDefault,
	}
	assert.NoError(t, repos.UserRepo.CreateUser(user))

	schema := &types.Schema{
		UID:     uuid.NewString(),
		UserUID: user.UID,
		SizeX:   10,
		SizeY:   10,
		SizeZ:   10,
	}
	assert.NoError(t, repos.SchemaRepo.CreateSchema(schema))

	m := core.NewSchemaMap(repos.SchemaPartRepo, schema)
	assert.NotNil(t, m)

	n, err := m.GetNode(mt.NewPos(0, 0, 0))
	assert.NoError(t, err)
	assert.Nil(t, n)

	assert.NoError(t, m.SetNode(mt.NewPos(1, 2, 3), &mt.Node{Name: "default:stone", Param1: 15, Param2: 2}, nil))

	n, err = m.GetNode(mt.NewPos(1, 2, 3))
	assert.NoError(t, err)
	assert.NotNil(t, n)
	assert.Equal(t, "default:stone", n.Name)
	assert.Equal(t, 15, n.Param1)
	assert.Equal(t, 2, n.Param2)

	assert.NoError(t, m.Close())
}
