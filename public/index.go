package public

import (
	"html/template"
	"net/http"
	"os"
)

type IndexModel struct {
	Webdev bool
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
		Webdev: os.Getenv("WEBDEV") == "true",
	}
	t.Execute(w, m)
}
