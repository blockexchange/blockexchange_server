package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestServe(t *testing.T) {
	r := mux.NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	api := createTestApi(t)
	SetupRoutes(r, api)

	res, err := http.Get(ts.URL + "/api/info")
	assert.NoError(t, err)

	info := Info{}
	err = json.NewDecoder(res.Body).Decode(&info)
	assert.NoError(t, err)
	assert.Equal(t, 1, info.VersionMajor)
	assert.Equal(t, 1, info.VersionMinor)
}
