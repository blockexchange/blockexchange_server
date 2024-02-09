package types

type SchemaMod struct {
	ID       *int64 `json:"id" ksql:"id"`
	SchemaID int64  `json:"schema_id" ksql:"schema_id"`
	ModName  string `json:"mod_name" ksql:"mod_name"`
}
