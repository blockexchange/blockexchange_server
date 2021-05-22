package web

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"blockexchange/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchSchema(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())

	user := testutils.CreateUser(api.UserRepo, t, &types.User{})
	schema := types.Schema{
		UserID: user.ID,
	}
	err := api.SchemaRepo.CreateSchema(&schema)
	assert.NoError(t, err)

	list, err := api.SchemaSearchRepo.FindByUsername(user.Name)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, schema.ID, list[0].ID)
}
