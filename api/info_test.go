package api

import (
	"blockexchange/core"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	i := InfoHandler{
		Config: &core.Config{},
	}
	i.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)

	info := Info{}
	err := json.NewDecoder(w.Body).Decode(&info)
	assert.NoError(t, err)
	assert.Equal(t, 1, info.VersionMajor)
	assert.Equal(t, 1, info.VersionMinor)
}
