package templateengine_test

import (
	"blockexchange/controller"
	"blockexchange/templateengine"
	"blockexchange/templateengine/testdata"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEngine(t *testing.T) {
	te := templateengine.NewTemplateEngine(&templateengine.TemplateEngineOptions{
		Templates:   testdata.Files,
		TemplateDir: "testdata",
		EnableCache: true,
		FuncMap:     make(map[string]any),
	})

	r := httptest.NewRequest("GET", "http://", nil)
	w := httptest.NewRecorder()

	assert.NoError(t, te.Execute("pages/error.html", w, r, 500, &controller.RenderData{BaseURL: "./", Data: errors.New("dummy")}))
	assert.NoError(t, te.Execute("pages/index.html", w, r, 200, &controller.RenderData{BaseURL: "./"}))
}
