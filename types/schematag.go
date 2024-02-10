package types

type SchemaTag struct {
	UID      string `json:"uid" ksql:"uid"`
	TagUID   string `json:"tag_uid" ksql:"tag_uid"`
	SchemaID int64  `json:"schema_id" ksql:"schema_id"`
}
