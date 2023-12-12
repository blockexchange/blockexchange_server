package api_test

import (
	"blockexchange/api"
	"blockexchange/testutils"
	"blockexchange/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewTestApi(t *testing.T) *api.Api {
	db_ := testutils.CreateTestDatabase(t)
	api, err := api.NewApi(db_, types.CreateConfig())
	assert.NoError(t, err)
	return api
}
