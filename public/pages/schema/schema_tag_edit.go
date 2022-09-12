package schema

import (
	"blockexchange/controller"
	"blockexchange/types"
	"errors"
	"fmt"
	"net/http"
)

type SchemaTagEditModel struct {
	Schema *types.SchemaSearchResult
	Tags   []*SchemaTagEditEntry
}

type SchemaTagEditEntry struct {
	*types.Tag
	Selected bool
}

func schemaTagEditPost(rc *controller.RenderContext) error {
	r := rc.Request()
	repos := rc.Repositories()
	claims := rc.Claims()

	r.ParseForm()
	schema, err := searchSchema(repos.SchemaSearchRepo, r)
	if err != nil {
		return err
	}
	if schema == nil {
		return errors.New("not found")
	}
	if claims.UserID != schema.UserID {
		return errors.New("not allowed")
	}

	schema_tags, err := repos.SchemaTagRepo.GetBySchemaID(schema.ID)
	if err != nil {
		return err
	}
	selectedTagIDS := make(map[int64]bool)
	for _, st := range schema_tags {
		selectedTagIDS[st.TagID] = true
	}

	tags, err := repos.TagRepo.GetAll()
	if err != nil {
		return err
	}

	for _, tag := range tags {
		label := fmt.Sprintf("tag_assigned_%d", tag.ID)
		if r.FormValue(label) != "" && !selectedTagIDS[tag.ID] {
			// new tag
			repos.SchemaTagRepo.Create(schema.ID, tag.ID)
		}
		if r.FormValue(label) == "" && selectedTagIDS[tag.ID] {
			// removed tag
			repos.SchemaTagRepo.Delete(schema.ID, tag.ID)
		}
	}

	rc.Redirect("../" + schema.Name)

	return nil
}

func SchemaTagEdit(rc *controller.RenderContext) error {
	r := rc.Request()

	if r.Method == http.MethodPost {
		return schemaTagEditPost(rc)
	}

	schema, err := searchSchema(rc.Repositories().SchemaSearchRepo, r)
	if err != nil {
		return err
	}
	m := &SchemaTagEditModel{Schema: schema}

	schema_tags, err := rc.Repositories().SchemaTagRepo.GetBySchemaID(schema.ID)
	if err != nil {
		return err
	}
	selectedTagIDS := make(map[int64]bool)
	for _, st := range schema_tags {
		selectedTagIDS[st.TagID] = true
	}

	tags, err := rc.Repositories().TagRepo.GetAll()
	if err != nil {
		return err
	}
	m.Tags = make([]*SchemaTagEditEntry, len(tags))
	for i, t := range tags {
		m.Tags[i] = &SchemaTagEditEntry{
			Tag:      t,
			Selected: selectedTagIDS[t.ID],
		}
	}

	return rc.Render("pages/schema/schema_tag_edit.html", m)
}
