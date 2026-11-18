package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestLegacyRedirects(t *testing.T) {
	r := mux.NewRouter()
	(&Api{}).SetupLegacyRedirects(r)

	tests := map[string]string{
		"/schema/someuser/my%20schema": "/#/schema/someuser/my%20schema",
		"/search?q=house&page=2":       "/#/search?q=house&page=2",
		"/profile":                     "/#/profile",
	}

	for path, expected := range tests {
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, path, nil))
		assert.Equal(t, http.StatusMovedPermanently, rr.Code, path)
		assert.Equal(t, expected, rr.Header().Get("Location"), path)
	}
}
