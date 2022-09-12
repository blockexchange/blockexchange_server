package api

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewTestApi(t *testing.T) *Api {
	db_ := testutils.CreateTestDatabase(t)
	api, err := NewApi(db_, core.NewNoOpCache())
	assert.NoError(t, err)
	return api
}
