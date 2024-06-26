package api

import (
	"blockexchange/core"
	"blockexchange/types"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// get user
func (api *Api) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_uid := vars["user_uid"]

	user, err := api.UserRepo.GetUserByUID(user_uid)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("get user error: %s", err.Error()))
		return
	}
	if user == nil {
		SendError(w, 404, fmt.Sprintf("user not found: %s", user_uid))
		return
	}

	Send(w, user, nil)
}

func (api *Api) CountUserSchemaStars(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_uid := vars["user_uid"]

	c, err := api.Repositories.SchemaStarRepo.CountByUserUID(user_uid)
	Send(w, c, err)
}

func (api *Api) CountUsers(w http.ResponseWriter, r *http.Request) {
	count, err := api.UserRepo.CountUsers()
	Send(w, count, err)
}

func (api *Api) SearchUsers(w http.ResponseWriter, r *http.Request) {
	search := &types.UserSearch{}
	err := json.NewDecoder(r.Body).Decode(search)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("search parse error: %s", err.Error()))
		return
	}

	result := []*types.User{}

	if search.Name != nil {
		// search by name
		user, err := api.UserRepo.GetUserByName(*search.Name)
		if err != nil {
			SendError(w, 500, fmt.Sprintf("search error: %s", err.Error()))
			return
		}
		if user != nil {
			result = append(result, user)
		}
	} else {
		// get all
		limit := 20
		offset := 0
		if search.Limit != nil && *search.Limit < 100 {
			limit = *search.Limit
		}
		if search.Offset != nil && *search.Offset >= 0 {
			offset = *search.Offset
		}

		result, err = api.UserRepo.GetUsers(limit, offset)
		if err != nil {
			SendError(w, 500, fmt.Sprintf("search-all error: %s", err.Error()))
			return
		}
	}

	Send(w, result, nil)
}

// save changed user (name, permissions)
func (api *Api) SaveUser(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	user_uid := vars["user_uid"]
	claims := ctx.Claims

	sent_user := &types.User{}
	err := json.NewDecoder(r.Body).Decode(sent_user)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("user decode error: %s", err.Error()))
		return
	}
	if sent_user.UID == "" {
		SendError(w, 500, "no id")
		return
	}

	// check if the requested user is the same
	if !claims.HasPermission(types.JWTPermissionAdmin) && sent_user.UID != claims.UserUID {
		SendError(w, 403, fmt.Sprintf("not authorized to get user '%s'", user_uid))
		return
	}

	user, err := api.UserRepo.GetUserByUID(user_uid)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("get user error: %s", err.Error()))
		return
	}
	if user == nil {
		SendError(w, 404, fmt.Sprintf("user not found: %s", user_uid))
		return
	}

	if claims.HasPermission(types.JWTPermissionAdmin) {
		// admin can save permissions
		user.Role = sent_user.Role
	}

	if sent_user.Name != user.Name {
		// name change, check conditions

		// valid name
		if !core.ValidateName(sent_user.Name) {
			SendError(w, 500, "invalid new name")
			return
		}

		// name available
		new_name_user, err := api.UserRepo.GetUserByName(sent_user.Name)
		if err != nil {
			SendError(w, 500, fmt.Sprintf("get user by new name error: %s", err.Error()))
			return
		}
		if new_name_user != nil {
			SendError(w, 500, fmt.Sprintf("user with name '%s' is already registered", sent_user.Name))
			return
		}

		user.Name = sent_user.Name
	}

	err = api.UserRepo.UpdateUser(user)
	Send(w, user, err)
}

func (api *Api) ChangePassword(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	user_uid := vars["user_uid"]

	chpwd := &types.ChangePassword{}
	err := json.NewDecoder(r.Body).Decode(chpwd)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("json parse error: %s", err.Error()))
		return
	}

	if len(chpwd.NewPassword) == 0 {
		SendError(w, 500, "new password empty")
		return
	}

	if !ctx.HasPermission(types.JWTPermissionAdmin) && ctx.Claims.UserUID != user_uid {
		SendError(w, 403, "unauthorized to change the password of another user")
		return
	}

	user, err := api.UserRepo.GetUserByUID(user_uid)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("get user error: %s", err.Error()))
		return
	}

	if user.Type != types.UserTypeLocal {
		SendError(w, 500, fmt.Sprintf("user type mismatch: %s", user.Type))
		return
	}

	if !ctx.HasPermission(types.JWTPermissionAdmin) {
		// check old password
		err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(chpwd.OldPassword))
		if err != nil {
			SendError(w, 401, err.Error())
			return
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(chpwd.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("hash error: %s", err.Error()))
		return
	}

	// save new password
	user.Hash = string(hash)
	err = api.UserRepo.UpdateUser(user)
	Send(w, true, err)
}

func (api *Api) UnlinkOauth(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	user_uid := vars["user_uid"]

	chpwd := &types.ChangePassword{}
	err := json.NewDecoder(r.Body).Decode(chpwd)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("json parse error: %s", err.Error()))
		return
	}

	if len(chpwd.NewPassword) == 0 {
		SendError(w, 500, "new password empty")
		return
	}

	if !ctx.HasPermission(types.JWTPermissionAdmin) && ctx.Claims.UserUID != user_uid {
		SendError(w, 403, "unauthorized to unlink another user")
		return
	}

	user, err := api.UserRepo.GetUserByUID(user_uid)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("get user error: %s", err.Error()))
		return
	}

	if user.Type == types.UserTypeLocal {
		SendError(w, 500, fmt.Sprintf("user type mismatch: %s", user.Type))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(chpwd.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("hash error: %s", err.Error()))
		return
	}

	// save user with new password
	user.Hash = string(hash)
	user.Type = types.UserTypeLocal
	err = api.UserRepo.UpdateUser(user)
	Send(w, true, err)
}
