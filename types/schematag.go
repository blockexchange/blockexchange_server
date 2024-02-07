package types

type SchemaTag struct {
	ID       int64 `json:"id" ksql:"id"`
	TagID    int64 `json:"tag_id" ksql:"tag_id"`
	SchemaID int64 `json:"schema_id" ksql:"schema_id"`
}
