package api

import (
	"blockexchange/types"
	"net/http"
	"strconv"
)

func GetLimitOffset(r *http.Request, default_limit int) (int, int) {
	q := r.URL.Query()
	if q.Get("limit") != "" && q.Get("offset") != "" {
		var limit, offset int64
		limit, err := strconv.ParseInt(q.Get("limit"), 10, 64)
		if err != nil {
			return default_limit, 0
		}

		offset, err = strconv.ParseInt(q.Get("offset"), 10, 64)
		if err != nil {
			return default_limit, 0
		}

		return int(limit), int(offset)
	} else {
		return default_limit, 0
	}
}

// loads the schema with the given uid and checks that the caller may modify it.
// Writes the error response and returns nil if the schema is missing (404),
// can't be loaded (500) or the caller is not the owner (403).
// If allow_admin is true, users with the admin permission pass the ownership check too.
func (api *Api) loadOwnedSchema(w http.ResponseWriter, ctx *SecureContext, schema_uid string, allow_admin bool) *types.Schema {
	schema, err := api.SchemaRepo.GetSchemaByUID(schema_uid)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	if schema == nil {
		SendError(w, http.StatusNotFound, "schema not found")
		return nil
	}

	is_owner := schema.UserUID == ctx.Claims.UserUID
	if !is_owner && !(allow_admin && ctx.HasPermission(types.JWTPermissionAdmin)) {
		SendError(w, http.StatusForbidden, "you are not the owner of the schema")
		return nil
	}

	return schema
}
