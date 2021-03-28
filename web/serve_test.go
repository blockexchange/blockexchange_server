package web

import (
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestServe(t *testing.T) {
	r := mux.NewRouter()
	ts := httptest.NewServer(r)

	api := createTestApi(t)
	SetupRoutes(r, api)

	//TODO: test stuff

	defer ts.Close()
}
