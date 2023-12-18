package api

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (api Api) UpdateSchemaTags(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if !ctx.CheckPermission(w, types.JWTPermissionManagement) {
		SendError(w, 403, "not a management token")
		return
	}

	schema, err := api.SchemaRepo.GetSchemaById(id)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if schema == nil {
		SendError(w, 404, "not found")
		return
	}

	// check permissions
	is_admin := ctx.CheckPermission(w, types.JWTPermissionAdmin)
	if !is_admin && schema.UserID != ctx.Claims.UserID {
		// not an admin and not the owner
		SendError(w, 403, "unauthorized")
		return
	}

	// assemble maps
	new_tag_names := []string{}
	err = json.NewDecoder(r.Body).Decode(&new_tag_names)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	new_tag_name_map := map[string]bool{}
	for _, new_tag_name := range new_tag_names {
		new_tag_name_map[new_tag_name] = true
	}

	tags, err := api.TagRepo.GetAll()
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	tag_id_name_map := map[int64]string{}
	tag_name_id_map := map[string]int64{}
	for _, t := range tags {
		tag_id_name_map[t.ID] = t.Name
		tag_name_id_map[t.Name] = t.ID
	}

	existing_tag_list, err := api.SchemaTagRepo.GetBySchemaID(schema.ID)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	existing_tags := map[string]bool{}
	for _, existing_tag := range existing_tag_list {
		name := tag_id_name_map[existing_tag.TagID]
		existing_tags[name] = true
	}

	// check for new tags
	for _, new_tag_name := range new_tag_names {
		if existing_tags[new_tag_name] {
			// still there
			continue
		}

		id := tag_name_id_map[new_tag_name]
		err = api.SchemaTagRepo.Create(schema.ID, id)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
	}

	// check for removed tags
	for existing_tag := range existing_tags {
		if !new_tag_name_map[existing_tag] {
			// tag removed
			id := tag_name_id_map[existing_tag]
			err = api.SchemaTagRepo.Delete(schema.ID, id)
			if err != nil {
				SendError(w, 500, err.Error())
				return
			}
		}
	}

	// send new list back
	Send(w, new_tag_names, nil)
}
