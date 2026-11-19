package api

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (api Api) UpdateSchemaTags(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	if !ctx.CheckPermission(w, types.JWTPermissionManagement) {
		return
	}

	schema := api.loadOwnedSchema(w, ctx, schema_uid, true)
	if schema == nil {
		return
	}
	is_admin := ctx.HasPermission(types.JWTPermissionAdmin)

	// assemble maps
	new_tag_names := []string{}
	err := json.NewDecoder(r.Body).Decode(&new_tag_names)
	if SendIfError(w, err) {
		return
	}
	new_tag_name_map := map[string]bool{}
	for _, new_tag_name := range new_tag_names {
		new_tag_name_map[new_tag_name] = true
	}

	tags, err := api.TagRepo.GetAll()
	if SendIfError(w, err) {
		return
	}
	tag_id_name_map := map[string]string{}
	tag_name_id_map := map[string]string{}
	restricted_id_map := map[string]bool{}
	for _, t := range tags {
		restricted_id_map[t.UID] = t.Restricted
		tag_id_name_map[t.UID] = t.Name
		tag_name_id_map[t.Name] = t.UID
	}

	existing_tag_list, err := api.SchemaTagRepo.GetBySchemaUID(schema.UID)
	if SendIfError(w, err) {
		return
	}
	existing_tags := map[string]bool{}
	for _, existing_tag := range existing_tag_list {
		name := tag_id_name_map[existing_tag.TagUID]
		existing_tags[name] = true
	}

	// check for new tags
	for _, new_tag_name := range new_tag_names {
		if existing_tags[new_tag_name] {
			// still there
			continue
		}

		id := tag_name_id_map[new_tag_name]
		if restricted_id_map[id] && !is_admin {
			// only admins can change restricted tags
			continue
		}

		err = api.SchemaTagRepo.Create(&types.SchemaTag{TagUID: id, SchemaUID: schema.UID})
		if SendIfError(w, err) {
			return
		}
	}

	// check for removed tags
	for existing_tag := range existing_tags {
		if !new_tag_name_map[existing_tag] {
			// tag removed
			id := tag_name_id_map[existing_tag]
			if restricted_id_map[id] && !is_admin {
				// only admins can change restricted tags
				continue
			}

			err = api.SchemaTagRepo.Delete(schema.UID, id)
			if SendIfError(w, err) {
				return
			}
		}
	}

	// send new list back
	Send(w, new_tag_names, nil)
}
