package api_test

import (
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodedefinitionsBatch(t *testing.T) {
	api := NewTestApi(t)
	assert.NoError(t, api.MediaRepo.UpsertNodedefinitions([]*types.Nodedefinition{
		{Name: "batchtest:stone", ModName: "batchtest", Definition: "{}"},
		{Name: "batchtest:dirt", ModName: "batchtest", Definition: "{}"},
	}))

	buf, err := json.Marshal([]string{"batchtest:stone", "batchtest:dirt", "unknown:node"})
	assert.NoError(t, err)

	r := httptest.NewRequest("POST", "http://", bytes.NewBuffer(buf))
	w := httptest.NewRecorder()
	api.GetNodedefinitionsBatch(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	list := []*types.Nodedefinition{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &list))
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "batchtest:dirt", list[0].Name)
	assert.Equal(t, "batchtest:stone", list[1].Name)

	// empty list
	r = httptest.NewRequest("POST", "http://", strings.NewReader("[]"))
	w = httptest.NewRecorder()
	api.GetNodedefinitionsBatch(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, "[]\n", w.Body.String())

	// invalid body
	r = httptest.NewRequest("POST", "http://", strings.NewReader("{"))
	w = httptest.NewRecorder()
	api.GetNodedefinitionsBatch(w, r)
	assert.Equal(t, 400, w.Result().StatusCode)

	// too many names
	buf, err = json.Marshal(make([]string, 1001))
	assert.NoError(t, err)
	r = httptest.NewRequest("POST", "http://", bytes.NewBuffer(buf))
	w = httptest.NewRecorder()
	api.GetNodedefinitionsBatch(w, r)
	assert.Equal(t, 400, w.Result().StatusCode)
}

func TestMediafilesBatch(t *testing.T) {
	api := NewTestApi(t)
	assert.NoError(t, api.MediaRepo.UpsertMediafiles([]*types.Mediafile{
		{Name: "batchtest_stone.png", ModName: "batchtest", Data: []byte{1, 2, 3}},
		{Name: "batchtest_dirt.png", ModName: "batchtest", Data: []byte{4, 5, 6}},
	}))

	buf, err := json.Marshal([]string{"batchtest_stone.png", "unknown.png"})
	assert.NoError(t, err)

	r := httptest.NewRequest("POST", "http://", bytes.NewBuffer(buf))
	w := httptest.NewRecorder()
	api.GetMediafilesBatch(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	list := []*types.Mediafile{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &list))
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "batchtest_stone.png", list[0].Name)
	assert.Equal(t, "batchtest", list[0].ModName)
	assert.Equal(t, []byte{1, 2, 3}, list[0].Data)
}
