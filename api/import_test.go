package api_test

import (
	"blockexchange/api"
	"blockexchange/testutils"
	"blockexchange/types"
	"bytes"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func importRequest(t *testing.T, a *api.Api, user *types.User, filename string, body []byte) int {
	r := httptest.NewRequest("POST", "http://", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{"filename": filename})

	a.ImportSchematic(w, r, &api.SecureContext{
		Claims: &types.Claims{Username: user.Name, UserUID: user.UID},
	})
	return w.Result().StatusCode
}

func TestImportWorldedit(t *testing.T) {
	a := NewTestApi(t)
	user := testutils.CreateUser(a.UserRepo, t, nil)

	data, err := os.ReadFile("../schematic/worldedit/testdata/plain_chest.we")
	assert.NoError(t, err)

	assert.Equal(t, 200, importRequest(t, a, user, "chest.we", data))
}

func TestImportInvalidWorldedit(t *testing.T) {
	a := NewTestApi(t)
	user := testutils.CreateUser(a.UserRepo, t, nil)

	assert.Equal(t, 400, importRequest(t, a, user, "x.we", []byte("5:while true do end")))
	assert.Equal(t, 400, importRequest(t, a, user, "x.we", []byte("5")))
	assert.Equal(t, 400, importRequest(t, a, user, "x.txt", []byte("hello")))
}

func TestImportTooLarge(t *testing.T) {
	a := NewTestApi(t)
	user := testutils.CreateUser(a.UserRepo, t, nil)

	assert.Equal(t, 413, importRequest(t, a, user, "x.we", make([]byte, api.MAX_SIZE+1)))
}

func TestImportBX(t *testing.T) {
	a := NewTestApi(t)
	user := testutils.CreateUser(a.UserRepo, t, nil)

	data, err := os.ReadFile("../schematic/bx/testdata/blockexchange.zip")
	assert.NoError(t, err)

	assert.Equal(t, 200, importRequest(t, a, user, "schema.zip", data))
	assert.Equal(t, 400, importRequest(t, a, user, "schema.zip", []byte("not a zip")))
}
