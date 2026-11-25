package api_test

import (
	"blockexchange/parser"
	"blockexchange/testutils"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCreateSchemaPartValid(t *testing.T) {
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, nil)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)
	schemapart := testutils.CreateSchemaPart(api.SchemaPartRepo, t, schema, nil)

	data, err := json.Marshal(schemapart)
	assert.NoError(t, err)

	// save

	r := httptest.NewRequest("GET", "http://", bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	Login(t, r, user)

	api.Secure(api.CreateSchemaPart)(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)

	// load

	r = httptest.NewRequest("GET", "http://", nil)
	w = httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"schema_uid": schemapart.SchemaUID,
		"x":          strconv.Itoa(int(schemapart.OffsetX)),
		"y":          strconv.Itoa(int(schemapart.OffsetY)),
		"z":          strconv.Itoa(int(schemapart.OffsetZ)),
	})

	api.GetSchemaPart(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	response_schemapart := &types.SchemaPart{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), response_schemapart))
	assert.Equal(t, schemapart.Data, response_schemapart.Data)
	assert.Equal(t, schemapart.MetaData, response_schemapart.MetaData)
}

func TestGetNextSchemaPart(t *testing.T) {
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, nil)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)
	testutils.CreateSchemaPart(api.SchemaPartRepo, t, schema, &types.SchemaPart{
		OffsetX:   0,
		OffsetY:   0,
		OffsetZ:   0,
		Mtime:     100,
		SchemaUID: schema.UID,
		Data:      []byte{},
		MetaData:  []byte{},
	})

	testutils.CreateSchemaPart(api.SchemaPartRepo, t, schema, &types.SchemaPart{
		OffsetX:   16,
		OffsetY:   0,
		OffsetZ:   0,
		Mtime:     200,
		SchemaUID: schema.UID,
		Data:      []byte{},
		MetaData:  []byte{},
	})

	// load first

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"schema_uid": schema.UID,
	})

	api.GetFirstSchemaPart(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	response_schemapart := &types.SchemaPart{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), response_schemapart))
	assert.Equal(t, 0, response_schemapart.OffsetX)
	assert.Equal(t, 0, response_schemapart.OffsetY)
	assert.Equal(t, 0, response_schemapart.OffsetZ)

	// load next

	r = httptest.NewRequest("GET", "http://", nil)
	w = httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"schema_uid": schema.UID,
		"x":          "0",
		"y":          "0",
		"z":          "0",
	})

	api.GetNextSchemaPart(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	response_schemapart = &types.SchemaPart{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), response_schemapart))
	assert.Equal(t, 16, response_schemapart.OffsetX)
	assert.Equal(t, 0, response_schemapart.OffsetY)
	assert.Equal(t, 0, response_schemapart.OffsetZ)

	// load by mtime

	r = httptest.NewRequest("GET", "http://", nil)
	w = httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"schema_uid": schema.UID,
		"mtime":      fmt.Sprintf("%d", 100),
	})

	api.GetNextSchemaPartByMtime(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	response_schemapart = &types.SchemaPart{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), response_schemapart))
	assert.Equal(t, int64(200), response_schemapart.Mtime)
	assert.Equal(t, 16, response_schemapart.OffsetX)
	assert.Equal(t, 0, response_schemapart.OffsetY)
	assert.Equal(t, 0, response_schemapart.OffsetZ)

	// count by mtime

	r = httptest.NewRequest("GET", "http://", nil)
	w = httptest.NewRecorder()

	r = mux.SetURLVars(r, map[string]string{
		"schema_uid": schema.UID,
		"mtime":      fmt.Sprintf("%d", 100),
	})

	api.CountNextSchemaPartByMtime(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, []byte{'1'}, w.Body.Bytes())
}

func TestCreateSchemaPartInvalidSchemaID(t *testing.T) {
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, nil)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)
	schemapart := testutils.CreateSchemaPart(api.SchemaPartRepo, t, schema, nil)
	schemapart.SchemaUID = ""

	data, err := json.Marshal(schemapart)
	assert.NoError(t, err)

	r := httptest.NewRequest("GET", "http://", bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	Login(t, r, user)

	api.Secure(api.CreateSchemaPart)(w, r)

	assert.Equal(t, 500, w.Result().StatusCode)
}

func TestGetSchemaPartChunk(t *testing.T) {
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, nil)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)

	// parts inside the first chunk (0..63) and at the start of the next one
	for _, x := range []int{0, 48, 64} {
		testutils.CreateSchemaPart(api.SchemaPartRepo, t, schema, &types.SchemaPart{
			OffsetX:   x,
			SchemaUID: schema.UID,
			Data:      []byte{},
			MetaData:  []byte{},
		})
	}

	get_chunk := func(x int) []*types.SchemaPart {
		r := httptest.NewRequest("GET", "http://", nil)
		w := httptest.NewRecorder()
		r = mux.SetURLVars(r, map[string]string{
			"schema_uid": schema.UID,
			"x":          strconv.Itoa(x),
			"y":          "0",
			"z":          "0",
		})
		api.GetSchemaPartChunk(w, r)
		assert.Equal(t, 200, w.Result().StatusCode)
		parts := []*types.SchemaPart{}
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &parts))
		return parts
	}

	offsets := func(parts []*types.SchemaPart) []int {
		list := []int{}
		for _, p := range parts {
			list = append(list, p.OffsetX)
		}
		return list
	}

	// the chunks don't overlap
	assert.ElementsMatch(t, []int{0, 48}, offsets(get_chunk(0)))
	assert.ElementsMatch(t, []int{64}, offsets(get_chunk(64)))
}

func TestPatchSchemaPart(t *testing.T) {
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, nil)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)

	newPart := func(name string, param2 byte) *types.SchemaPart {
		psp := &parser.ParsedSchemaPart{
			NodeIDS:        []int16{1, 1},
			Param1:         []byte{0, 0},
			Param2:         []byte{param2, param2},
			NodeNameLookup: map[int16]string{1: name},
			Meta: &parser.SchemaPartMetadata{
				NodeMapping: map[string]int16{name: 1},
				Size:        parser.SchemaPartSize{X: 2, Y: 1, Z: 1},
			},
		}
		sp, err := psp.Convert()
		assert.NoError(t, err)
		sp.SchemaUID = schema.UID
		sp.Mtime = 100
		return sp
	}

	testutils.CreateSchemaPart(api.SchemaPartRepo, t, schema, newPart("default:stone", 1))

	patch := &types.PatchedSchemapart{
		SchemaPart:  newPart("default:dirt", 2),
		PatchFilter: []types.PatchFilterType{types.PatchReplaceAll, types.PatchReplaceParam2},
	}
	data, err := json.Marshal(patch)
	assert.NoError(t, err)

	r := httptest.NewRequest("PATCH", "http://", bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	Login(t, r, user)
	api.Secure(api.PatchSchemaPart)(w, r)
	assert.Equal(t, 200, w.Result().StatusCode)

	response_schemapart := &types.SchemaPart{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), response_schemapart))
	assert.Greater(t, response_schemapart.Mtime, int64(100))

	stored, err := api.SchemaPartRepo.GetBySchemaUIDAndOffset(schema.UID, 0, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, response_schemapart.Data, stored.Data)
	assert.Equal(t, response_schemapart.Mtime, stored.Mtime)

	parsed, err := parser.ParseSchemaPart(stored)
	assert.NoError(t, err)
	assert.Equal(t, "default:dirt", parsed.NodeNameLookup[parsed.NodeIDS[0]])
	assert.Equal(t, "default:stone", parsed.NodeNameLookup[parsed.NodeIDS[1]])
	assert.Equal(t, []byte{2, 2}, parsed.Param2)

	// other user is not allowed to patch
	other := testutils.CreateUser(api.UserRepo, t, nil)
	r = httptest.NewRequest("PATCH", "http://", bytes.NewBuffer(data))
	w = httptest.NewRecorder()
	Login(t, r, other)
	api.Secure(api.PatchSchemaPart)(w, r)
	assert.Equal(t, 403, w.Result().StatusCode)

	// missing part
	patch.OffsetX = 16
	data, err = json.Marshal(patch)
	assert.NoError(t, err)
	r = httptest.NewRequest("PATCH", "http://", bytes.NewBuffer(data))
	w = httptest.NewRecorder()
	Login(t, r, user)
	api.Secure(api.PatchSchemaPart)(w, r)
	assert.Equal(t, 404, w.Result().StatusCode)
}
