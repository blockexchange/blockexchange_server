package types

var OrderDirections = map[string]bool{
	"asc":  true,
	"desc": true,
}

var OrderColumns = map[string]bool{
	"s.created": true,
	"s.mtime":   true,
	"s.stars":   true,
}

type SchemaSearchRequest struct {
	UserID         *int64  `json:"user_id"`
	SchemaID       *int64  `json:"schema_id"`
	TagID          *int64  `json:"tag_id"`
	SchemaName     *string `json:"schema_name"`
	UserName       *string `json:"user_name"`
	Keywords       *string `json:"keywords"`
	Complete       *bool   `json:"complete"`
	OrderDirection *string `json:"order_direction"`
	OrderColumn    *string `json:"order_column"`
	Limit          *int    `json:"limit"`
	Offset         *int    `json:"offset"`
}
