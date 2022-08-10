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
		JWTKey:       "mykey",
		CookieName:   "my-app",
		CookiePath:   "/",
		CookieDomain: "127.0.0.1",
		CookieSecure: false,
	})

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	assert.NoError(t, te.ExecuteError(w, r, "./", 500, errors.New("dummy")))
	assert.NoError(t, te.Execute("pages/index.html", w, r, "./", nil))

}
