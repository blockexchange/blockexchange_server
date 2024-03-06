package types

type SchemaMod struct {
	SchemaUID string `json:"schema_uid" ksql:"schema_uid"`
	ModName   string `json:"mod_name" ksql:"mod_name"`
}

type SchemaModCount struct {
	ModName string `json:"mod_name" ksql:"mod_name"`
	Count   int64  `json:"count" ksql:"count"`
}
