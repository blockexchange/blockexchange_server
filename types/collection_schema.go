package types

type CollectionSchema struct {
	CollectionID int64 `json:"collection_id" db:"collection_id"`
	SchemaID     int64 `json:"schema_id" db:"schema_id"`
}
