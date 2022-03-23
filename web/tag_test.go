package web

import (
	"blockexchange/testutils"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTags(t *testing.T) {
	api := NewTestApi(t)

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

func TestCreateTag(t *testing.T) {
	api := NewTestApi(t)

	tag := types.Tag{
		Name:        "test",
		Description: "123",
	}

	data, err := json.Marshal(tag)
	assert.NoError(t, err)
	r := httptest.NewRequest("GET", "http://", bytes.NewReader(data))
	w := httptest.NewRecorder()
	testutils.Login(t, r, &types.User{
		Name: "admin",
		Role: types.UserRoleAdmin,
	})

	Secure(api.CreateTag)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)
	fmt.Println(w.Result())
}
