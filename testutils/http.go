package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func CreateGETRequest(t *testing.T, vars map[string]string) (*httptest.ResponseRecorder, *http.Request) {
	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	if vars != nil {
		r = mux.SetURLVars(r, vars)
	}
	return w, r
}

func CreatePOSTRequest(t *testing.T, vars map[string]string, o any) (*httptest.ResponseRecorder, *http.Request) {

	var buf *bytes.Buffer = nil

	if o != nil {
		data, err := json.Marshal(o)
		assert.NoError(t, err)
		buf = bytes.NewBuffer(data)
	}

	r := httptest.NewRequest("POST", "http://", buf)
	w := httptest.NewRecorder()

	if vars != nil {
		r = mux.SetURLVars(r, vars)
	}
	return w, r
}

func CreateDELETERequest(t *testing.T, vars map[string]string) (*httptest.ResponseRecorder, *http.Request) {
	r := httptest.NewRequest("DELETE", "http://", nil)
	w := httptest.NewRecorder()

	if vars != nil {
		r = mux.SetURLVars(r, vars)
	}
	return w, r
}
