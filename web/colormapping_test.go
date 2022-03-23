package web

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetColorMapping(t *testing.T) {
	api := NewTestApi(t)

	// Create
	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	api.GetColorMapping(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
}
