package web

import (
	"blockexchange/core"
	"blockexchange/testutils"
	"blockexchange/types"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTags(t *testing.T) {
	db_ := testutils.CreateTestDatabase(t)
	api := NewApi(db_, core.NewNoOpCache())

	err := api.TagRepo.Create(&types.Tag{
		Name:        "test",
		Description: "123",
	})
	assert.NoError(t, err)

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	api.GetTags(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)

	tags := []types.Tag{}
	err = json.NewDecoder(w.Body).Decode(&tags)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, "test", tags[0].Name)
	assert.Equal(t, "123", tags[0].Description)
}
