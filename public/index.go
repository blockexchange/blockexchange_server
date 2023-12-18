package public

import (
	"html/template"
	"net/http"
	"os"
)

type IndexModel struct {
	Webdev  bool
	BaseURL string
	Meta    map[string]string
}

func RenderIndex(w http.ResponseWriter, r *http.Request, meta map[string]string) {
	data, err := Webapp.ReadFile("index.html")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	t, err := template.New("").Parse(string(data))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	m := &IndexModel{
		Webdev:  os.Getenv("WEBDEV") == "true",
		BaseURL: os.Getenv("BASE_URL"),
		Meta:    meta,
	}
	t.Execute(w, m)
}
