package web

import (
	"blockexchange/types"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const template_str = `
<!DOCTYPE HTML>
<html>
	<head>
		<meta name="og:title" content="{{.Schema.Name}} by {{.Username}}"/>
		<meta name="og:type" content="Schematic"/>
		<meta name="og:url" content="{{.BaseURL}}/#/schema/{{.Username}}/{{.Schema.Name}}"/>
		<meta name="og:image" content="{{.BaseURL}}/api/schema/{{.Schema.ID}}/screenshot/{{.Screenshot.ID}}"/>
		<meta name="og:site_name" content="Block exchange"/>
		<meta name="og:description" content="{{.Schema.Description}}"/>
		<meta http-equiv="refresh" content="0; url={{.BaseURL}}/#/schema/{{.Username}}/{{.Schema.Name}}" />
	</head>
	<body>
	</body>
</html>
`

var tmpl = template.Must(template.New("main").Parse(template_str))

type TemplateData struct {
	Schema     *types.SchemaSearchResult
	Screenshot *types.SchemaScreenshot
	Username   string
	BaseURL    string
}

func (api *Api) GetStaticView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_name := vars["schema_name"]
	user_name := vars["user_name"]

	schema, err := api.SchemaSearchRepo.FindByUsernameAndSchemaname(schema_name, user_name)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if schema == nil {
		SendError(w, 404, "Not found")
		return
	}

	screenshots, err := api.SchemaScreenshotRepo.GetBySchemaID(schema.ID)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if screenshots == nil || len(screenshots) < 1 {
		SendError(w, 500, "no screenshots found")
		return
	}

	data := TemplateData{
		Schema:     schema,
		Screenshot: &screenshots[0],
		Username:   user_name,
		BaseURL:    os.Getenv("BASE_URL"),
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}
