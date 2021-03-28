package types

type SchemaTag struct {
	ID       int64 `json:"id" db:"id"`
	TagID    int64 `json:"tag_id" db:"tag_id"`
	SchemaID int64 `json:"schema_id" db:"schema_id"`
}
