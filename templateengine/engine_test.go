package templateengine_test

import (
	"blockexchange/public"
	"blockexchange/templateengine"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEngine(t *testing.T) {
	te := templateengine.NewTemplateEngine(&templateengine.TemplateEngineOptions{
		Templates:    public.Files,
		TemplateDir:  "public",
		EnableCache:  true,
		CookieName:   "my-app",
		CookiePath:   "/",
		CookieDomain: "127.0.0.1",
		CookieSecure: false,
		FuncMap:      make(map[string]any),
	})

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	assert.NoError(t, te.ExecuteError(w, r, "./", 500, errors.New("dummy")))
}
