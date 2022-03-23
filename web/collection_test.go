package web

import (
	"blockexchange/testutils"
	"blockexchange/types"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollection(t *testing.T) {
	api := NewTestApi(t)

	// empty collection

	user := testutils.CreateUser(api.UserRepo, t, nil)
	w, r := testutils.CreateGETRequest(t, map[string]string{
		"user_id": fmt.Sprintf("%d", user.ID),
	})

	api.GetCollectionsByUserID(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	list := []types.Collection{}
	assert.NoError(t, json.NewDecoder(w.Body).Decode(&list))
	assert.Equal(t, 0, len(list))

	// create collection

	c := &types.Collection{
		Name:        "mycollection",
		Description: "something, something",
	}
	w, r = testutils.CreatePOSTRequest(t, nil, c)
	testutils.Login(t, r, user)

	Secure(api.CreateCollection)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	// TODO: update

	// populated collection

	w, r = testutils.CreateGETRequest(t, map[string]string{
		"user_id": fmt.Sprintf("%d", user.ID),
	})

	api.GetCollectionsByUserID(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.NoError(t, json.NewDecoder(w.Body).Decode(&list))
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "mycollection", c.Name)
	assert.Equal(t, "something, something", c.Description)
	created_collection_id := list[0].ID

	// delete collection

	w, r = testutils.CreateDELETERequest(t, map[string]string{
		"id": fmt.Sprintf("%d", created_collection_id),
	})
	testutils.Login(t, r, user)
	Secure(api.DeleteCollection)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	// empty collection

	w, r = testutils.CreateGETRequest(t, map[string]string{
		"user_id": fmt.Sprintf("%d", user.ID),
	})

	api.GetCollectionsByUserID(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	list = []types.Collection{}
	assert.NoError(t, json.NewDecoder(w.Body).Decode(&list))
	assert.Equal(t, 0, len(list))

}
