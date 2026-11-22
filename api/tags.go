package api

import (
	"blockexchange/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// editable fields of a tag
type TagRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Restricted  bool   `json:"restricted"`
}

func (req *TagRequest) validate() error {
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(req.Name) > 128 {
		return fmt.Errorf("name must be at most 128 characters")
	}
	return nil
}

func (api Api) GetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := api.TagRepo.GetAll()
	Send(w, tags, err)
}

// checks that no other tag than the one with the given uid has the name
func (api Api) checkTagNameAvailable(w http.ResponseWriter, name string, uid string) bool {
	existing, err := api.TagRepo.GetByName(name)
	if SendIfError(w, err) {
		return false
	}
	if existing != nil && existing.UID != uid {
		SendError(w, http.StatusConflict, fmt.Sprintf("a tag named '%s' already exists", name))
		return false
	}
	return true
}

func decodeTagRequest(w http.ResponseWriter, r *http.Request) *TagRequest {
	req := &TagRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		SendError(w, http.StatusBadRequest, fmt.Sprintf("json error: %s", err.Error()))
		return nil
	}
	if err := req.validate(); err != nil {
		SendError(w, http.StatusBadRequest, err.Error())
		return nil
	}
	return req
}

func (api Api) CreateTag(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionAdmin) {
		return
	}
	req := decodeTagRequest(w, r)
	if req == nil {
		return
	}
	if !api.checkTagNameAvailable(w, req.Name, "") {
		return
	}

	tag := &types.Tag{
		Name:        req.Name,
		Description: req.Description,
		Restricted:  req.Restricted,
	}
	err := api.TagRepo.Create(tag)
	Send(w, tag, err)
}

func (api Api) UpdateTag(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionAdmin) {
		return
	}
	tag, err := api.TagRepo.GetByUID(mux.Vars(r)["tag_uid"])
	if SendIfError(w, err) {
		return
	}
	if tag == nil {
		SendError(w, http.StatusNotFound, "tag not found")
		return
	}
	req := decodeTagRequest(w, r)
	if req == nil {
		return
	}
	if !api.checkTagNameAvailable(w, req.Name, tag.UID) {
		return
	}

	tag.Name = req.Name
	tag.Description = req.Description
	tag.Restricted = req.Restricted
	err = api.TagRepo.Update(tag)
	Send(w, tag, err)
}
