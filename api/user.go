package api

import (
	"blockexchange/core"
	"blockexchange/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// get user
func (api *Api) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_id := vars["user_id"]

	id, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("user_id parse error: %s", err.Error()))
		return
	}

	user, err := api.UserRepo.GetUserById(id)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("get user error: %s", err.Error()))
		return
	}
	if user == nil {
		SendError(w, 404, fmt.Sprintf("user not found: %d", id))
		return
	}

	Send(w, user, nil)
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
	user_id := vars["user_id"]
	claims := ctx.Claims

	sent_user := &types.User{}
	err := json.NewDecoder(r.Body).Decode(sent_user)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("user decode error: %s", err.Error()))
		return
	}
	if sent_user.ID == nil {
		SendError(w, 500, "no id")
		return
	}

	// check if the requested user is the same
	if !claims.HasPermission(types.JWTPermissionAdmin) && *sent_user.ID != claims.UserID {
		SendError(w, 403, fmt.Sprintf("not authorized to get user '%s'", user_id))
		return
	}

	id, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("user_id parse error: %s", err.Error()))
		return
	}

	user, err := api.UserRepo.GetUserById(id)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("get user error: %s", err.Error()))
		return
	}
	if user == nil {
		SendError(w, 404, fmt.Sprintf("user not found: %d", id))
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
