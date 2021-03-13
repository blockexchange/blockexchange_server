package types

type SchemaMod struct {
	ID       int64  `json:"id" db:"id"`
	SchemaID int64  `json:"schema_id" db:"schema_id"`
	ModName  string `json:"mod_name" db:"mod_name"`
}
