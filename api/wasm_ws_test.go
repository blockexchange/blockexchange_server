package api_test

import (
	"blockexchange/api"
	"blockexchange/core"
	"blockexchange/testutils"
	"blockexchange/types"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWasmWebsocket(t *testing.T) {
	cfg := types.CreateConfig()
	cfg.BaseURL = "http://127.0.0.1:1"
	a, router, err := api.NewApi(testutils.CreateTestDatabase(t), cfg)
	assert.NoError(t, err)

	owner := testutils.CreateUser(a.UserRepo, t, nil)
	other := testutils.CreateUser(a.UserRepo, t, nil)
	schema := testutils.CreateSchema(a.SchemaRepo, t, owner, nil)

	// the full router, including the logging/metrics middleware
	srv := httptest.NewServer(router)
	defer srv.Close()
	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/api/schema/" + schema.UID + "/wasm/ws"

	c := core.New(cfg, nil)
	dial := func(user *types.User, perms ...types.JWTPermission) (*websocket.Conn, int) {
		h := http.Header{}
		if user != nil {
			token, err := c.CreateJWT(user, perms, time.Hour)
			assert.NoError(t, err)
			h.Set("Authorization", token)
		}
		ws, resp, err := websocket.DefaultDialer.Dial(url, h)
		if err != nil {
			return nil, resp.StatusCode
		}
		return ws, resp.StatusCode
	}

	_, status := dial(nil)
	assert.Equal(t, 401, status)
	_, status = dial(other)
	assert.Equal(t, 403, status)

	ws, status := dial(other, types.JWTPermissionAdmin)
	assert.Equal(t, 101, status)
	ws.Close()

	ws, status = dial(owner)
	assert.Equal(t, 101, status)
	defer ws.Close()
	ws.SetReadDeadline(time.Now().Add(5 * time.Second))
	assert.NoError(t, ws.WriteMessage(websocket.TextMessage, []byte("PROXY IPV4 TCP 10.0.0.1 8080")))
	_, data, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, "PROXY OK", string(data))

	// cross-origin upgrades are refused
	h := http.Header{}
	token, _ := c.CreateJWT(owner, nil, time.Hour)
	h.Set("Authorization", token)
	h.Set("Origin", "https://evil.example.com")
	_, resp, err := websocket.DefaultDialer.Dial(url, h)
	assert.Error(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}
