package public

import (
	"html/template"
	"net/http"
	"os"
)

type IndexModel struct {
	Webdev  bool
	BaseURL string
}

func RenderIndex(w http.ResponseWriter, r *http.Request) {
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
	}
	t.Execute(w, m)
}
