package types

type CollectionSearchResponse struct {
	*Collection
	Username string   `json:"username"`
	Size     int      `json:"size"`
	Tags     []string `json:"tags"`
	Mods     []string `json:"mods"`
}

type CollectionSearchRequest struct {
	UserUID   *string `json:"user_uid"`
	SchemaUID *string `json:"schema_uid"`
	TagUID    *string `json:"tag_uid"`
	Keywords  *string `json:"keywords"`
	Limit     *int    `json:"limit"`
	Offset    *int    `json:"offset"`
}
