package api_test

import (
	"blockexchange/api"
	"blockexchange/schematic/mtschem"
	"blockexchange/testutils"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestWorldEditExport(t *testing.T) {
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, nil)
	assert.NotEmpty(t, user.UID)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)
	testutils.CreateSchemaPart(api.SchemaPartRepo, t, schema, nil)

	w, r := testutils.CreateGETRequest(t, map[string]string{
		"schema_uid": schema.UID,
	})

	api.ExportWorldeditSchema(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", w.Result().Header.Get("Content-Type"))
}

func TestBXExport(t *testing.T) {
	api := NewTestApi(t)
	user := testutils.CreateUser(api.UserRepo, t, nil)
	schema := testutils.CreateSchema(api.SchemaRepo, t, user, nil)
	testutils.CreateSchemaPart(api.SchemaPartRepo, t, schema, nil)

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{
		"schema_uid": schema.UID,
	})

	api.ExportBXSchema(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, "application/zip", w.Result().Header.Get("Content-Type"))
}

func TestMTSExportImport(t *testing.T) {
	a := NewTestApi(t)
	user := testutils.CreateUser(a.UserRepo, t, nil)

	data, err := os.ReadFile("../schematic/mtschem/testdata/St.Peters_Basilica(Vatican,Italy).mts")
	assert.NoError(t, err)

	// import
	r := mux.SetURLVars(httptest.NewRequest("POST", "http://", bytes.NewReader(data)), map[string]string{"filename": "basilica.mts"})
	w := httptest.NewRecorder()
	a.ImportSchematic(w, r, &api.SecureContext{Claims: &types.Claims{Username: user.Name, UserUID: user.UID}})
	assert.Equal(t, 200, w.Result().StatusCode, w.Body.String())
	schema := &types.Schema{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), schema))

	// export
	r = mux.SetURLVars(httptest.NewRequest("GET", "http://", nil), map[string]string{"schema_uid": schema.UID, "filename": "basilica.mts"})
	w = httptest.NewRecorder()
	a.ExportMTSSchema(w, r)
	assert.Equal(t, 200, w.Result().StatusCode, w.Body.String())
	assert.Equal(t, "application/octet-stream", w.Result().Header.Get("Content-Type"))

	original, err := mtschem.Parse(data)
	assert.NoError(t, err)
	exported, err := mtschem.Parse(w.Body.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, original.SizeX, exported.SizeX)
	assert.Equal(t, original.SizeY, exported.SizeY)
	assert.Equal(t, original.SizeZ, exported.SizeZ)

	// unknown schema
	r = mux.SetURLVars(httptest.NewRequest("GET", "http://", nil), map[string]string{"schema_uid": uuid.NewString(), "filename": "x.mts"})
	w = httptest.NewRecorder()
	a.ExportMTSSchema(w, r)
	assert.Equal(t, 404, w.Result().StatusCode)
}
