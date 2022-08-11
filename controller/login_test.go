package controller_test

import (
	"blockexchange/controller"
	"blockexchange/core"
	"blockexchange/testutils"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	db := testutils.CreateTestDatabase(t)
	cfg := &core.Config{}

	c := controller.NewController(db, cfg)
	r := httptest.NewRequest("POST", "http://", nil)
	w := httptest.NewRecorder()

	c.Login(w, r)

	assert.Equal(t, 200, w.Result().StatusCode)
}
