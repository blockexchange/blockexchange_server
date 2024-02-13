package types

type SchemaStar struct {
	UserUID  string `json:"user_uid" ksql:"user_uid"`
	SchemaID int64  `json:"schema_id" ksql:"schema_id"`
}

type SchemaStarResponse struct {
	Count   int  `json:"count"`
	Starred bool `json:"starred"`
}
