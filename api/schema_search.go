package api

import (
	"blockexchange/types"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (api *Api) AddSchemaSearchFields(schemas []*types.Schema) ([]*types.SchemaSearchResponse, error) {
	user_uids := []string{}
	schema_uids := []string{}

	for _, s := range schemas {
		user_uids = append(user_uids, s.UserUID)
		schema_uids = append(schema_uids, s.UID)
	}

	users, err := api.Repositories.UserRepo.GetUsersByUIDs(user_uids)
	if err != nil {
		return nil, err
	}
	user_map := map[string]*types.User{}
	for _, u := range users {
		user_map[u.UID] = u
	}

	list := make([]*types.SchemaSearchResponse, len(schemas))
	schema_map := map[string]*types.SchemaSearchResponse{}
	for i, s := range schemas {
		user := user_map[s.UserUID]
		if user == nil {
			return nil, fmt.Errorf("user-id %s not found", s.UserUID)
		}
		sr := &types.SchemaSearchResponse{
			Schema:   s,
			Username: user.Name,
			Tags:     []string{},
			Mods:     []string{},
		}
		schema_map[s.UID] = sr
		list[i] = sr
	}

	tags, err := api.TagRepo.GetAll()
	if err != nil {
		return nil, err
	}
	tag_map := map[string]*types.Tag{}
	for _, t := range tags {
		tag_map[t.UID] = t
	}

	schema_tags, err := api.Repositories.SchemaTagRepo.GetBySchemaUIDs(schema_uids)
	if err != nil {
		return nil, err
	}
	for _, st := range schema_tags {
		sr := schema_map[st.SchemaUID]
		if sr == nil {
			return nil, fmt.Errorf("schema %s for schema-tag %s not found", st.SchemaUID, st.UID)
		}
		t := tag_map[st.TagUID]
		if t == nil {
			return nil, fmt.Errorf("tag %s for schema-tag %s not found", st.TagUID, st.UID)
		}
		sr.Tags = append(sr.Tags, t.Name)
	}

	schema_mods, err := api.Repositories.SchemaModRepo.GetSchemaModsBySchemaUIDs(schema_uids)
	if err != nil {
		return nil, err
	}
	for _, sm := range schema_mods {
		sr := schema_map[sm.SchemaUID]
		if sr == nil {
			return nil, fmt.Errorf("schema %s for schema-mod not found", sm.SchemaUID)
		}
		sr.Mods = append(sr.Mods, sm.ModName)
	}

	return list, nil
}

func (api *Api) SearchSchemaByNameAndUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_name := vars["schema_name"]
	user_name := vars["user_name"]
	limit := 1
	offset := 0

	search := &types.SchemaSearchRequest{
		UserName:   &user_name,
		SchemaName: &schema_name,
		Limit:      &limit,
		Offset:     &offset,
	}
	list, err := api.SchemaSearchRepo.Search(search)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if len(list) == 0 {
		SendError(w, 404, "not found")
		return
	}

	list2, err := api.AddSchemaSearchFields(list)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	schema := list2[0]
	if r.URL.Query().Get("download") == "true" {
		// increment downloads and ignore error
		api.incrementDownloadstats(schema.UID, r)
	}

	Send(w, schema, nil)
}

func (api *Api) CountSchema(w http.ResponseWriter, r *http.Request) {
	search := &types.SchemaSearchRequest{}
	err := json.NewDecoder(r.Body).Decode(search)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	c, err := api.SchemaSearchRepo.Count(search)
	Send(w, c, err)
}

func (api *Api) SearchSchema(w http.ResponseWriter, r *http.Request) {
	search := &types.SchemaSearchRequest{}
	err := json.NewDecoder(r.Body).Decode(search)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// apply sane defaults
	if search.Limit == nil || *search.Limit > 100 || *search.Limit <= 0 {
		l := 100
		search.Limit = &l
	}
	if search.Offset == nil || *search.Offset > 10000 || *search.Offset < 0 {
		o := 0
		search.Offset = &o
	}

	list, err := api.SchemaSearchRepo.Search(search)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	list2, err := api.AddSchemaSearchFields(list)
	Send(w, list2, err)
}
