package types

type SchemaStar struct {
	UserUID   string `json:"user_uid" ksql:"user_uid"`
	SchemaUID string `json:"schema_uid" ksql:"schema_uid"`
}

type SchemaStarResponse struct {
	Count   int  `json:"count"`
	Starred bool `json:"starred"`
}
