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
	Schema         *Schema  `json:"schema"`
	CollectionName *string  `json:"collection_name"`
	Username       string   `json:"username"`
	Tags           []string `json:"tags"`
	Mods           []string `json:"mods"`
}

type SchemaSearchRequest struct {
	UserUID        *string `json:"user_uid"`
	SchemaUID      *string `json:"schema_uid"`
	TagUID         *string `json:"tag_uid"`
	TagName        *string `json:"tag_name"`
	SchemaName     *string `json:"schema_name"`
	ModName        *string `json:"mod_name"`
	UserName       *string `json:"user_name"`
	Keywords       *string `json:"keywords"`
	Complete       *bool   `json:"complete"`
	CollectionUID  *string `json:"collection_uid"`
	CollectionName *string `json:"collection_name"`
	WithCollection *bool   `json:"with_collection"`
	OrderDirection *string `json:"order_direction"`
	OrderColumn    *string `json:"order_column"`
	FromMtime      *int64  `json:"from_mtime"`
	UntilMtime     *int64  `json:"until_mtime"`
	Limit          *int    `json:"limit"`
	Offset         *int    `json:"offset"`
}
