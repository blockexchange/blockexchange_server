package types

var OrderDirections = map[string]bool{
	"asc":  true,
	"desc": true,
}

var OrderColumns = map[string]bool{
	"created": true,
	"mtime":   true,
	"stars":   true,
}

type SchemaSearchResponse struct {
	*Schema
	Username string   `json:"username"`
	Tags     []string `json:"tags"`
	Mods     []string `json:"mods"`
}

type SchemaSearchRequest struct {
	UserUID        *string `json:"user_uid"`
	SchemaID       *int64  `json:"schema_id"`
	TagUID         *string `json:"tag_uid"`
	SchemaName     *string `json:"schema_name"`
	UserName       *string `json:"user_name"`
	Keywords       *string `json:"keywords"`
	Complete       *bool   `json:"complete"`
	OrderDirection *string `json:"order_direction"`
	OrderColumn    *string `json:"order_column"`
	Limit          *int    `json:"limit"`
	Offset         *int    `json:"offset"`
}
