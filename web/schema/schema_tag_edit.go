package schema

import (
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

func (sc *SchemaContext) schemaTagEditPost(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	r.ParseForm()
	schema, err := searchSchema(sc.repos.SchemaSearchRepo, r)
	if err != nil {
		sc.tu.RenderError(w, r, 500, err)
		return
	}
	if schema == nil {
		sc.tu.RenderError(w, r, 404, errors.New("not found"))
		return
	}
	if claims.UserID != schema.UserID {
		sc.tu.RenderError(w, r, 403, errors.New("not allowed"))
		return
	}

	schema_tags, err := sc.repos.SchemaTagRepo.GetBySchemaID(schema.ID)
	if err != nil {
		sc.tu.RenderError(w, r, 500, err)
		return
	}
	selectedTagIDS := make(map[int64]bool)
	for _, st := range schema_tags {
		selectedTagIDS[st.TagID] = true
	}

	tags, err := sc.repos.TagRepo.GetAll()
	if err != nil {
		sc.tu.RenderError(w, r, 500, err)
		return
	}

	for _, tag := range tags {
		label := fmt.Sprintf("tag_assigned_%d", tag.ID)
		if r.FormValue(label) != "" && !selectedTagIDS[tag.ID] {
			// new tag
			sc.repos.SchemaTagRepo.Create(schema.ID, tag.ID)
		}
		if r.FormValue(label) == "" && selectedTagIDS[tag.ID] {
			// removed tag
			sc.repos.SchemaTagRepo.Delete(schema.ID, tag.ID)
		}
	}

	http.Redirect(w, r, sc.BaseURL+"/schema/"+schema.UserName+"/"+schema.Name, http.StatusSeeOther)
}

func (sc *SchemaContext) SchemaTagEdit(w http.ResponseWriter, r *http.Request, claims *types.Claims) {

	if r.Method == http.MethodPost {
		sc.schemaTagEditPost(w, r, claims)
		return
	}

	schema, err := searchSchema(sc.repos.SchemaSearchRepo, r)
	if err != nil {
		sc.tu.RenderError(w, r, 500, err)
		return
	}
	m := &SchemaTagEditModel{Schema: schema}

	schema_tags, err := sc.repos.SchemaTagRepo.GetBySchemaID(schema.ID)
	if err != nil {
		sc.tu.RenderError(w, r, 500, err)
		return
	}
	selectedTagIDS := make(map[int64]bool)
	for _, st := range schema_tags {
		selectedTagIDS[st.TagID] = true
	}

	tags, err := sc.repos.TagRepo.GetAll()
	if err != nil {
		sc.tu.RenderError(w, r, 500, err)
		return
	}
	m.Tags = make([]*SchemaTagEditEntry, len(tags))
	for i, t := range tags {
		m.Tags[i] = &SchemaTagEditEntry{
			Tag:      t,
			Selected: selectedTagIDS[t.ID],
		}
	}

	sc.tu.ExecuteTemplate(w, r, "schema/schema_tag_edit.html", m)
}
