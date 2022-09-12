package core

import (
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

const template_str = `
Schema created: **{{.Schema.Name}}** by **{{.User.Name}}**
Link: {{.BaseURL}}/schema/{{.User.Name}}/{{.Schema.Name}}
`

var feed_template = template.Must(template.New("main").Parse(template_str))

type TemplateData struct {
	Schema     *types.Schema
	User       *types.User
	BaseURL    string
	CodeMarker string
}

type DiscordData struct {
	Content string `json:"content"`
}

func renderFeedTemplate(baseUrl string, schema *types.Schema, user *types.User, screenshot *types.SchemaScreenshot) (*bytes.Buffer, error) {
	data := TemplateData{
		Schema:     schema,
		User:       user,
		BaseURL:    baseUrl,
		CodeMarker: "```",
	}

	buf := bytes.NewBuffer([]byte{})
	err := feed_template.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// posts the new schema to a discord channel
// NOTE: errors here are only logged not escalated
func UpdateSchemaFeed(schema *types.Schema, user *types.User, screenshot *types.SchemaScreenshot) {
	feed_url := os.Getenv("DISCORD_SCHEMA_FEED_URL")
	if feed_url == "" {
		// not configured
		return
	}

	baseUrl := os.Getenv("BASE_URL")
	buf, err := renderFeedTemplate(baseUrl, schema, user, screenshot)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"schema_id": schema.ID,
			"err":       err.Error(),
		}).Error("UpdateSchemaFeed::Template")
		return
	}

	discord_data := DiscordData{
		Content: buf.String(),
	}

	json_data, err := json.Marshal(&discord_data)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"schema_id": schema.ID,
			"err":       err.Error(),
		}).Error("UpdateSchemaFeed::Marshal")
		return
	}

	req, err := http.NewRequest("POST", feed_url, bytes.NewBuffer(json_data))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"schema_id": schema.ID,
			"err":       err.Error(),
		}).Error("UpdateSchemaFeed::NewRequest")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"schema_id": schema.ID,
			"err":       err.Error(),
		}).Error("UpdateSchemaFeed::ExecRequest")
		return
	}

	if resp.StatusCode != 200 {
		logrus.WithFields(logrus.Fields{
			"schema_id": schema.ID,
			"status":    resp.StatusCode,
		}).Error("UpdateSchemaFeed::Response")
		return
	}

}
