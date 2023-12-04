package api

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"blockexchange/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewTestApi(t *testing.T) *Api {
	db_ := testutils.CreateTestDatabase(t)
	api, err := NewApi(db_, core.NewNoOpCache(), types.CreateConfig())
	assert.NoError(t, err)
	return api
}
