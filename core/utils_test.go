package core_test

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/testutils"
	"blockexchange/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestCore(t *testing.T) (*core.Core, *db.Repositories) {
	kdb := testutils.CreateTestDatabase(t)
	repos := db.NewRepositories(kdb)

	c := core.New(types.CreateConfig(), repos)
	assert.NotNil(t, c)
	return c, repos
}
