package api_test

import (
	"blockexchange/api"
	"blockexchange/testutils"
	"blockexchange/types"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWasmRuntimeEmbedded(t *testing.T) {
	cfg := types.CreateConfig()
	_, router, err := api.NewApi(testutils.CreateTestDatabase(t), cfg)
	assert.NoError(t, err)

	get := func(path string, accept_encoding string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", path, nil)
		if accept_encoding != "" {
			r.Header.Set("Accept-Encoding", accept_encoding)
		}
		router.ServeHTTP(w, r)
		return w
	}

	// stored as luanti.wasm.gz, served compressed
	w := get("/wasm/runtime/luanti.wasm", "gzip, br")
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
	assert.Equal(t, "application/wasm", w.Header().Get("Content-Type"))
	assert.Equal(t, "require-corp", w.Header().Get("Cross-Origin-Embedder-Policy"))
	assert.Less(t, w.Body.Len(), 6*1024*1024)

	// decompressed for clients without gzip support
	w = get("/wasm/runtime/luanti.wasm", "")
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "", w.Header().Get("Content-Encoding"))
	assert.Equal(t, []byte{0, 'a', 's', 'm'}, w.Body.Bytes()[:4])

	w = get("/wasm/runtime/luanti.js", "gzip")
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "javascript")

	for _, p := range []string{"/wasm/runtime/worker.js", "/wasm/runtime/packs/base.pack", "/wasm/runtime/packs/certs.pack"} {
		assert.Equal(t, 200, get(p, "gzip").Code, p)
	}
	assert.Equal(t, 404, get("/wasm/runtime/", "").Code)
	assert.Equal(t, 404, get("/wasm/runtime/packs/", "").Code)
	assert.Equal(t, 404, get("/wasm/runtime/missing.js", "").Code)

	// the webapp as cross-origin isolated document for the "Edit Online" page, the normal one is not isolated
	w = get("/wasm/", "")
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `<div id="app">`)
	assert.Equal(t, "require-corp", w.Header().Get("Cross-Origin-Embedder-Policy"))
	assert.Equal(t, "same-origin", w.Header().Get("Cross-Origin-Opener-Policy"))
	w = get("/", "")
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `<div id="app">`)
	assert.Equal(t, "", w.Header().Get("Cross-Origin-Embedder-Policy"))

	w = get("/api/info", "")
	info := map[string]any{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &info))
	assert.Equal(t, true, info["wasm_enabled"])
}
