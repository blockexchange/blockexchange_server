package web

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTags(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	api.GetTags(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
}
