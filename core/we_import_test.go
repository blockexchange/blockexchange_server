package core_test

import (
	"blockexchange/testutils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWEImport(t *testing.T) {
	c, repos := getTestCore(t)

	u := testutils.CreateUser(repos.UserRepo, t, nil)
	assert.NotNil(t, u)

	data, err := os.ReadFile("testdata/plain_chest.we")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	s, err := c.ImportWE(data, u.Name, "we_import")
	assert.NoError(t, err)
	assert.NotNil(t, s)
}
