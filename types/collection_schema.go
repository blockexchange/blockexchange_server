package types

type CollectionSchema struct {
	CollectionID int64 `json:"collection_id"`
	SchemaID     int64 `json:"schema_id"`
}

func (cs *CollectionSchema) Table() string {
	return "collection_schema"
}

func (cs *CollectionSchema) Columns(action string) []string {
	return []string{"collection_id", "schema_id"}
}

func (cs *CollectionSchema) Values(action string) []any {
	return []any{cs.CollectionID, cs.SchemaID}
}

func (cs *CollectionSchema) Scan(action string, r func(dest ...any) error) error {
	return r(&cs.CollectionID, &cs.SchemaID)
}
