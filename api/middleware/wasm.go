package middleware

import (
	"net/http"
	"strings"
)

const (
	// webapp index (cross-origin isolated copy for the "Edit Online" page)
	WasmPrefix = "/wasm/"
	// luanti wasm build (embedded, see wasm/README.md)
	WasmRuntimePrefix = "/wasm/runtime/"
)

// the luanti wasm client needs SharedArrayBuffer, which is only available in a cross-origin isolated page.
// The headers are limited to the wasm paths, the rest of the webapp loads cross-origin images (avatars, previews)
func WasmHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, WasmPrefix) {
			h := w.Header()
			h.Set("Cross-Origin-Opener-Policy", "same-origin")
			h.Set("Cross-Origin-Embedder-Policy", "require-corp")
			h.Set("Cross-Origin-Resource-Policy", "same-origin")
			// always revalidate, the runtime and launcher have no versioned file names
			h.Set("Cache-Control", "no-cache")
		}
		next.ServeHTTP(w, r)
	})
}
