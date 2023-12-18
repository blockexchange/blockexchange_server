package core

import (
	"blockexchange/db"
	"blockexchange/testutils"
	"blockexchange/types"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWEImport(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	repos := db.NewRepositories(db_)

	c := New(types.CreateConfig(), repos)
	assert.NotNil(t, c)

	u := testutils.CreateUser(repos.UserRepo, t, nil)
	assert.NotNil(t, u)

	data, err := os.ReadFile("testdata/plain_chest.we")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	s, err := c.ImportWE(data, u.Name)
	assert.NoError(t, err)
	assert.NotNil(t, s)
}
