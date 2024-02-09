package types

type SchemaStar struct {
	UserID   int64 `json:"user_id" ksql:"user_id"`
	SchemaID int64 `json:"schema_id" ksql:"schema_id"`
}

type SchemaStarResponse struct {
	Count   int  `json:"count"`
	Starred bool `json:"starred"`
}
