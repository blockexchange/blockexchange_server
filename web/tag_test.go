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

	tag := types.Tag{
		Name:        "test",
		Description: "123",
	}
	err := api.TagRepo.Create(&tag)
	assert.NoError(t, err)

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	api.GetTags(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)

	list, err := api.TagRepo.GetAll()
	assert.NoError(t, err)

	tags := []types.Tag{}
	err = json.NewDecoder(w.Body).Decode(&tags)
	assert.NoError(t, err)
	assert.Equal(t, len(list), len(tags))

	err = api.TagRepo.Delete(tag.ID)
	assert.NoError(t, err)
}
