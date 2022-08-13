package templateengine_test

import (
	"blockexchange/templateengine"
	"blockexchange/templateengine/testdata"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEngine(t *testing.T) {
	te := templateengine.NewTemplateEngine(&templateengine.TemplateEngineOptions{
		Templates:    testdata.Files,
		TemplateDir:  "testdata",
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
	assert.NoError(t, te.Execute("pages/index.html", w, r, "./", nil))
}
