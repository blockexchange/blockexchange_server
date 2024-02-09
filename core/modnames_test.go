package core_test

import (
	"blockexchange/testutils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractModnames(t *testing.T) {
	c, repos := getTestCore(t)

	u := testutils.CreateUser(repos.UserRepo, t, nil)
	assert.NotNil(t, u)

	data, err := os.ReadFile("testdata/blockexchange.zip")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	s, err := c.ImportBX(data, u.Name)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	nodenames, err := c.ExtractModnames(*s.ID)
	assert.NoError(t, err)
	assert.NotNil(t, nodenames)
	assert.Equal(t, 1, len(nodenames))
	assert.Contains(t, nodenames, "default")
}
