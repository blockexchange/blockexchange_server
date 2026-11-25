package core_test

import (
	"blockexchange/core"
	"blockexchange/parser"
	"blockexchange/testutils"
	"blockexchange/types"
	"testing"

	mt "github.com/luanti-golang/types"

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

	// untouched nodes of a new part are air
	n, err = m.GetNode(mt.NewPos(0, 0, 0))
	assert.NoError(t, err)
	assert.Equal(t, "air", n.Name)

	assert.NoError(t, m.Close())

	// reload from the database
	m = core.NewSchemaMap(repos.SchemaPartRepo, schema)
	n, err = m.GetNode(mt.NewPos(0, 0, 0))
	assert.NoError(t, err)
	assert.Equal(t, "air", n.Name)
	n, err = m.GetNode(mt.NewPos(1, 2, 3))
	assert.NoError(t, err)
	assert.Equal(t, "default:stone", n.Name)
}

func TestSchemaMapSetNodeExistingMapping(t *testing.T) {
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
	md := &parser.MetadataEntry{Fields: map[string]string{"infotext": "chest"}}

	assert.NoError(t, m.SetNode(mt.NewPos(0, 0, 0), &mt.Node{Name: "default:stone"}, nil))
	assert.NoError(t, m.SetNode(mt.NewPos(1, 0, 0), &mt.Node{Name: "default:dirt"}, nil))
	// name already in mapping, metadata on a position without previous entry
	assert.NoError(t, m.SetNode(mt.NewPos(2, 0, 0), &mt.Node{Name: "default:stone"}, md))
	assert.NoError(t, m.Close())

	// reload from the database
	m = core.NewSchemaMap(repos.SchemaPartRepo, schema)
	for pos, name := range map[*mt.Pos]string{
		mt.NewPos(0, 0, 0): "default:stone",
		mt.NewPos(1, 0, 0): "default:dirt",
		mt.NewPos(2, 0, 0): "default:stone",
	} {
		n, err := m.GetNode(pos)
		assert.NoError(t, err)
		assert.NotNil(t, n)
		assert.Equal(t, name, n.Name, pos)
	}

	// overwrite an existing node of a stored part with an already mapped name
	assert.NoError(t, m.SetNode(mt.NewPos(1, 0, 0), &mt.Node{Name: "default:stone"}, md))
	n, err := m.GetNode(mt.NewPos(1, 0, 0))
	assert.NoError(t, err)
	assert.Equal(t, "default:stone", n.Name)
	assert.NoError(t, m.Close())

	part, err := repos.SchemaPartRepo.GetBySchemaUIDAndOffset(schema.UID, 0, 0, 0)
	assert.NoError(t, err)
	psp, err := parser.ParseSchemaPart(part)
	assert.NoError(t, err)
	assert.Equal(t, "chest", psp.Meta.Metadata.Meta["(1,0,0)"].Fields["infotext"])
	assert.Equal(t, "chest", psp.Meta.Metadata.Meta["(2,0,0)"].Fields["infotext"])
}
