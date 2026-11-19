package api

import (
	"blockexchange/core"
	"blockexchange/types"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// editable fields of a luanti media connection, the last_updated_* fields are set by the media sync
type LuantiMediaConnectionRequest struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req *LuantiMediaConnectionRequest) validate() error {
	if req.Hostname == "" {
		return fmt.Errorf("hostname is required")
	}
	if req.Port < 1 || req.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if req.Username == "" {
		return fmt.Errorf("username is required")
	}
	return nil
}

// the stored password is never sent to the client
func redactLuantiMediaConnection(c *types.LuantiMediaConnection) {
	c.Password = ""
}

func (api Api) GetLuantiMediaConnections(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionMedia) {
		return
	}
	list, err := api.LuantiMediaConnectionRepo.GetAll()
	if SendIfError(w, err) {
		return
	}
	for _, c := range list {
		redactLuantiMediaConnection(c)
	}
	Send(w, list, nil)
}

func (api Api) CreateLuantiMediaConnection(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionMedia) {
		return
	}
	req := &LuantiMediaConnectionRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		SendError(w, http.StatusBadRequest, fmt.Sprintf("json error: %s", err.Error()))
		return
	}
	if err := req.validate(); err != nil {
		SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	c := &types.LuantiMediaConnection{
		Hostname: req.Hostname,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
	}
	if SendIfError(w, api.LuantiMediaConnectionRepo.Create(c)) {
		return
	}
	redactLuantiMediaConnection(c)
	Send(w, c, nil)
}

func (api Api) UpdateLuantiMediaConnection(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionMedia) {
		return
	}
	c, err := api.LuantiMediaConnectionRepo.GetByID(mux.Vars(r)["id"])
	if SendIfError(w, err) {
		return
	}
	if c == nil {
		SendError(w, http.StatusNotFound, "connection not found")
		return
	}

	req := &LuantiMediaConnectionRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		SendError(w, http.StatusBadRequest, fmt.Sprintf("json error: %s", err.Error()))
		return
	}
	if err := req.validate(); err != nil {
		SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	c.Hostname = req.Hostname
	c.Port = req.Port
	c.Username = req.Username
	if req.Password != "" {
		// an empty password keeps the stored one
		c.Password = req.Password
	}
	if SendIfError(w, api.LuantiMediaConnectionRepo.Update(c)) {
		return
	}
	redactLuantiMediaConnection(c)
	Send(w, c, nil)
}

func (api Api) DeleteLuantiMediaConnection(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionMedia) {
		return
	}
	err := api.LuantiMediaConnectionRepo.Delete(mux.Vars(r)["id"])
	Send(w, true, err)
}

// fetches the node definitions and media files of the connection right away and returns the updated status
func (api Api) RunLuantiMediaConnection(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	if !ctx.CheckPermission(w, types.JWTPermissionMedia) {
		return
	}
	c, err := api.LuantiMediaConnectionRepo.GetByID(mux.Vars(r)["id"])
	if SendIfError(w, err) {
		return
	}
	if c == nil {
		SendError(w, http.StatusNotFound, "connection not found")
		return
	}

	// a failed sync is reported in the last_updated_* fields, not as an error response
	err = api.core.SyncLuantiMediaConnection(c)
	if errors.Is(err, core.ErrLuantiSyncRunning) {
		SendError(w, http.StatusConflict, err.Error())
		return
	}
	if SendIfError(w, err) {
		return
	}
	redactLuantiMediaConnection(c)
	Send(w, c, nil)
}
