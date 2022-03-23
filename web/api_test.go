package web

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"testing"
)

func NewTestApi(t *testing.T) *Api {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())
	return api
}
