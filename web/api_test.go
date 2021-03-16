// +build integration

package web

import (
	"blockexchange/db"
	"blockexchange/types"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
	//go test ./... -tags=integration
	_db, err := db.Init()
	assert.NoError(t, err)
	api := NewApi(_db)

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	api.GetUsers(w, r)
	var users []types.User
	err = json.NewDecoder(w.Body).Decode(&users)
	assert.NoError(t, err)
}
