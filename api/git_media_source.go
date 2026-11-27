package api

import (
	"blockexchange/core"
	"blockexchange/types"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/gorilla/mux"
)

// editable fields of a git media source, the last_updated_* fields are set by the sync,
// empty license fields are fetched from the contentdb on the next sync
type GitMediaSourceRequest struct {
	GitURL       string `json:"git_url"`
	CDBPackage   string `json:"cdb_package"`
	MediaLicense string `json:"media_license"`
	License      string `json:"license"`
}

var cdbPackageRegex = regexp.MustCompile(`^[A-Za-z0-9_.\-]+/[A-Za-z0-9_.\-]+$`)

func (req *GitMediaSourceRequest) validate() error {
	u, err := url.Parse(req.GitURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return fmt.Errorf("git_url must be a http(s) url")
	}
	if len(req.GitURL) > 256 {
		return fmt.Errorf("git_url is too long, max 256 characters")
	}
	if req.CDBPackage != "" && (len(req.CDBPackage) > 128 || !cdbPackageRegex.MatchString(req.CDBPackage)) {
		return fmt.Errorf("cdb_package must be in the form 'author/name'")
	}
	if len(req.MediaLicense) > 64 || len(req.License) > 64 {
		return fmt.Errorf("license names are limited to 64 characters")
	}
	return nil
}

func (req *GitMediaSourceRequest) apply(s *types.GitMediaSource) {
	s.GitURL = req.GitURL
	s.CDBPackage = req.CDBPackage
	s.MediaLicense = req.MediaLicense
	s.License = req.License
}

func decodeGitMediaSourceRequest(w http.ResponseWriter, r *http.Request) (*GitMediaSourceRequest, bool) {
	req := &GitMediaSourceRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		SendError(w, http.StatusBadRequest, fmt.Sprintf("json error: %s", err.Error()))
		return nil, false
	}
	if err := req.validate(); err != nil {
		SendError(w, http.StatusBadRequest, err.Error())
		return nil, false
	}
	return req, true
}

// the git url is unique
func (api Api) checkGitMediaSourceURL(w http.ResponseWriter, git_url, id string) bool {
	list, err := api.GitMediaSourceRepo.GetAll()
	if SendIfError(w, err) {
		return false
	}
	for _, s := range list {
		if s.GitURL == git_url && s.ID != id {
			SendError(w, http.StatusConflict, "a source with this git url already exists")
			return false
		}
	}
	return true
}

func (api Api) GetGitMediaSources(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionMedia) {
		return
	}
	list, err := api.GitMediaSourceRepo.GetAll()
	Send(w, list, err)
}

func (api Api) CreateGitMediaSource(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionMedia) {
		return
	}
	req, ok := decodeGitMediaSourceRequest(w, r)
	if !ok || !api.checkGitMediaSourceURL(w, req.GitURL, "") {
		return
	}
	s := &types.GitMediaSource{}
	req.apply(s)
	if SendIfError(w, api.GitMediaSourceRepo.Create(s)) {
		return
	}
	Send(w, s, nil)
}

func (api Api) UpdateGitMediaSource(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionMedia) {
		return
	}
	s, err := api.GitMediaSourceRepo.GetByID(mux.Vars(r)["id"])
	if SendIfError(w, err) {
		return
	}
	if s == nil {
		SendError(w, http.StatusNotFound, "source not found")
		return
	}
	req, ok := decodeGitMediaSourceRequest(w, r)
	if !ok || !api.checkGitMediaSourceURL(w, req.GitURL, s.ID) {
		return
	}
	req.apply(s)
	if SendIfError(w, api.GitMediaSourceRepo.Update(s)) {
		return
	}
	Send(w, s, nil)
}

// removes the source and its media files
func (api Api) DeleteGitMediaSource(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionMedia) {
		return
	}
	err := api.GitMediaSourceRepo.Delete(mux.Vars(r)["id"])
	Send(w, true, err)
}

// imports the media files of the source right away and returns the updated status
func (api Api) RunGitMediaSource(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionMedia) {
		return
	}
	s, err := api.GitMediaSourceRepo.GetByID(mux.Vars(r)["id"])
	if SendIfError(w, err) {
		return
	}
	if s == nil {
		SendError(w, http.StatusNotFound, "source not found")
		return
	}

	// a failed sync is reported in the last_updated_* fields, not as an error response
	err = api.core.SyncGitMediaSource(s)
	if errors.Is(err, core.ErrGitMediaSyncRunning) {
		SendError(w, http.StatusConflict, err.Error())
		return
	}
	Send(w, s, err)
}
