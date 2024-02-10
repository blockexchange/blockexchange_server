package types

type SchemaMod struct {
	UID      string `json:"uid" ksql:"uid"`
	SchemaID int64  `json:"schema_id" ksql:"schema_id"`
	ModName  string `json:"mod_name" ksql:"mod_name"`
}
