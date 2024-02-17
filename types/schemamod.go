package types

type SchemaMod struct {
	SchemaUID string `json:"schema_uid" ksql:"schema_uid"`
	ModName   string `json:"mod_name" ksql:"mod_name"`
}
